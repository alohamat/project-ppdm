package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"ppdm/models"
)

type DealsService struct {
	httpClient *http.Client
	cacheMutex sync.RWMutex
	cachedData *models.DealsResponse
	lastFetch  time.Time
	cacheTTL   time.Duration
}

func NewDealsService() *DealsService {
	return &DealsService{
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
		},
		cacheTTL: 10 * time.Minute,
	}
}

// Structs internas para deserialização da API da Epic Games
type epicResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []epicElement `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

type epicElement struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductSlug string `json:"productSlug"`
	UrlSlug     string `json:"urlSlug"`
	OfferType   string `json:"offerType"`
	KeyImages   []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"keyImages"`
	CatalogNs struct {
		Mappings []struct {
			PageSlug string `json:"pageSlug"`
			PageType string `json:"pageType"`
		} `json:"mappings"`
	} `json:"catalogNs"`
	OfferMappings []struct {
		PageSlug string `json:"pageSlug"`
		PageType string `json:"pageType"`
	} `json:"offerMappings"`
	Price struct {
		TotalPrice struct {
			DiscountPrice int    `json:"discountPrice"`
			OriginalPrice int    `json:"originalPrice"`
			CurrencyCode  string `json:"currencyCode"`
		} `json:"totalPrice"`
	} `json:"price"`
	Promotions *struct {
		PromotionalOffers []struct {
			PromotionalOffers []struct {
				StartDate       time.Time `json:"startDate"`
				EndDate         time.Time `json:"endDate"`
				DiscountSetting struct {
					DiscountType       string `json:"discountType"`
					DiscountPercentage int    `json:"discountPercentage"`
				} `json:"discountSetting"`
			} `json:"promotionalOffers"`
		} `json:"promotionalOffers"`
		UpcomingPromotionalOffers []struct {
			PromotionalOffers []struct {
				StartDate       time.Time `json:"startDate"`
				EndDate         time.Time `json:"endDate"`
				DiscountSetting struct {
					DiscountType       string `json:"discountType"`
					DiscountPercentage int    `json:"discountPercentage"`
				} `json:"discountSetting"`
			} `json:"promotionalOffers"`
		} `json:"upcomingPromotionalOffers"`
	} `json:"promotions"`
}

// Structs internas para deserialização da API da Steam
type steamFeaturedCategoriesResponse struct {
	Specials struct {
		ID    string          `json:"id"`
		Name  string          `json:"name"`
		Items []steamDealItem `json:"items"`
	} `json:"specials"`
	TopSellers struct {
		Items []steamDealItem `json:"items"`
	} `json:"top_sellers"`
}

type steamDealItem struct {
	ID                 int    `json:"id"`
	Type               int    `json:"type"`
	Name               string `json:"name"`
	Discounted         bool   `json:"discounted"`
	DiscountPercent    int    `json:"discount_percent"`
	OriginalPrice      int    `json:"original_price"`
	FinalPrice         int    `json:"final_price"`
	Currency           string `json:"currency"`
	LargeCapsuleImage  string `json:"large_capsule_image"`
	SmallCapsuleImage  string `json:"small_capsule_image"`
	HeaderImage        string `json:"header_image"`
	DiscountExpiration int64  `json:"discount_expiration"`
}

func (s *DealsService) GetDeals(ctx context.Context) (*models.DealsResponse, error) {
	s.cacheMutex.RLock()
	if s.cachedData != nil && time.Since(s.lastFetch) < s.cacheTTL {
		defer s.cacheMutex.RUnlock()
		return s.cachedData, nil
	}
	s.cacheMutex.RUnlock()

	return s.refreshDeals(ctx)
}

func (s *DealsService) refreshDeals(ctx context.Context) (*models.DealsResponse, error) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Double check após adquirir o lock de escrita
	if s.cachedData != nil && time.Since(s.lastFetch) < s.cacheTTL {
		return s.cachedData, nil
	}

	freeNow := make([]models.GameDeal, 0)
	freeUpcoming := make([]models.GameDeal, 0)
	bigDeals := make([]models.GameDeal, 0)

	// Busca paralela em Epic e Steam
	var wg sync.WaitGroup
	var epicErr, steamErr error
	var epicDeals, epicUpcoming, epicSpecials []models.GameDeal
	var steamSpecials, steamFree []models.GameDeal

	wg.Add(2)

	go func() {
		defer wg.Done()
		epicDeals, epicUpcoming, epicSpecials, epicErr = s.fetchEpicGames(ctx)
		if epicErr != nil {
			log.Printf("[DealsService] Erro ao buscar Epic Games: %v\n", epicErr)
		}
	}()

	go func() {
		defer wg.Done()
		steamSpecials, steamFree, steamErr = s.fetchSteamSpecials(ctx)
		if steamErr != nil {
			log.Printf("[DealsService] Erro ao buscar Steam: %v\n", steamErr)
		}
	}()

	wg.Wait()

	if epicErr != nil && steamErr != nil && s.cachedData != nil {
		// Se ambos falharem mas temos cache anterior, mantemos o cache para tolerância a falhas
		log.Println("[DealsService] Ambas as fontes falharam, utilizando cache anterior expirado.")
		return s.cachedData, nil
	}

	freeNow = append(freeNow, epicDeals...)
	freeNow = append(freeNow, steamFree...)

	freeUpcoming = append(freeUpcoming, epicUpcoming...)

	bigDeals = append(bigDeals, steamSpecials...)
	bigDeals = append(bigDeals, epicSpecials...)

	result := &models.DealsResponse{
		FreeNow:      freeNow,
		FreeUpcoming: freeUpcoming,
		BigDeals:     bigDeals,
		UpdatedAt:    time.Now(),
	}

	s.cachedData = result
	s.lastFetch = time.Now()

	return result, nil
}

func (s *DealsService) fetchEpicGames(ctx context.Context) (freeNow []models.GameDeal, upcoming []models.GameDeal, specials []models.GameDeal, err error) {
	url := "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=pt-BR&country=BR&allowCountries=BR"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao criar requisição epic: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao chamar api epic games: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("epic games retornou status %d", resp.StatusCode)
	}

	var epic epicResponse
	if err := json.NewDecoder(resp.Body).Decode(&epic); err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao decodificar resposta da epic games: %w", err)
	}

	now := time.Now()

	for _, elem := range epic.Data.Catalog.SearchStore.Elements {
		// Seleciona melhor imagem disponível
		var imageURL string
		for _, img := range elem.KeyImages {
			if strings.HasPrefix(img.URL, "http") {
				if img.Type == "OfferImageWide" || img.Type == "DieselStoreFrontWide" {
					imageURL = img.URL
					break
				}
				if imageURL == "" {
					imageURL = img.URL
				}
			}
		}

		// Determina o slug para o link na loja
		var slug string
		if len(elem.OfferMappings) > 0 && elem.OfferMappings[0].PageSlug != "" {
			slug = elem.OfferMappings[0].PageSlug
		} else if len(elem.CatalogNs.Mappings) > 0 && elem.CatalogNs.Mappings[0].PageSlug != "" {
			slug = elem.CatalogNs.Mappings[0].PageSlug
		} else if elem.ProductSlug != "" {
			slug = elem.ProductSlug
		} else if elem.UrlSlug != "" {
			slug = elem.UrlSlug
		}

		storeURL := "https://store.epicgames.com/pt-BR"
		if slug != "" {
			cleanSlug := strings.TrimSuffix(slug, "/home")
			storeURL = fmt.Sprintf("https://store.epicgames.com/pt-BR/p/%s", cleanSlug)
		}

		origPrice := float64(elem.Price.TotalPrice.OriginalPrice) / 100.0
		discPrice := float64(elem.Price.TotalPrice.DiscountPrice) / 100.0

		// Checa se está grátis agora via promoções
		isFreeNow := false
		var promoStart, promoEnd *time.Time

		if elem.Promotions != nil {
			for _, group := range elem.Promotions.PromotionalOffers {
				for _, offer := range group.PromotionalOffers {
					if offer.DiscountSetting.DiscountPercentage == 0 || discPrice == 0 {
						if (offer.StartDate.Before(now) || offer.StartDate.Equal(now)) && offer.EndDate.After(now) {
							isFreeNow = true
							st := offer.StartDate
							en := offer.EndDate
							promoStart = &st
							promoEnd = &en
							break
						}
					}
				}
				if isFreeNow {
					break
				}
			}
		}

		// Checa se é grátis diretamente pelo preço promocional
		if !isFreeNow && discPrice == 0 && origPrice > 0 {
			isFreeNow = true
		}

		if isFreeNow {
			freeNow = append(freeNow, models.GameDeal{
				ID:              elem.ID,
				Title:           elem.Title,
				Store:           "epic",
				Status:          "free_now",
				OriginalPrice:   origPrice,
				FinalPrice:      0,
				DiscountPercent: 100,
				Currency:        "BRL",
				ImageUrl:        imageURL,
				StoreUrl:        storeURL,
				StartDate:       promoStart,
				EndDate:         promoEnd,
			})
			continue
		}

		// Checa se será grátis no futuro próximo (upcoming)
		isUpcomingFree := false
		var upStart, upEnd *time.Time

		if elem.Promotions != nil {
			for _, group := range elem.Promotions.UpcomingPromotionalOffers {
				for _, offer := range group.PromotionalOffers {
					if offer.DiscountSetting.DiscountPercentage == 0 {
						if offer.StartDate.After(now) {
							isUpcomingFree = true
							st := offer.StartDate
							en := offer.EndDate
							upStart = &st
							upEnd = &en
							break
						}
					}
				}
				if isUpcomingFree {
					break
				}
			}
		}

		if isUpcomingFree {
			upcoming = append(upcoming, models.GameDeal{
				ID:              elem.ID,
				Title:           elem.Title,
				Store:           "epic",
				Status:          "free_upcoming",
				OriginalPrice:   origPrice,
				FinalPrice:      0,
				DiscountPercent: 100,
				Currency:        "BRL",
				ImageUrl:        imageURL,
				StoreUrl:        storeURL,
				StartDate:       upStart,
				EndDate:         upEnd,
			})
			continue
		}

		// Checa grandes promoções (deals) da Epic Games
		if origPrice > 0 && discPrice > 0 && discPrice < origPrice {
			pct := int(math.Round(((origPrice - discPrice) / origPrice) * 100.0))
			if pct >= 20 { // Apenas ofertas com desconto relevante
				specials = append(specials, models.GameDeal{
					ID:              elem.ID,
					Title:           elem.Title,
					Store:           "epic",
					Status:          "deal",
					OriginalPrice:   origPrice,
					FinalPrice:      discPrice,
					DiscountPercent: pct,
					Currency:        "BRL",
					ImageUrl:        imageURL,
					StoreUrl:        storeURL,
				})
			}
		}
	}

	return freeNow, upcoming, specials, nil
}

func (s *DealsService) fetchSteamSpecials(ctx context.Context) (specials []models.GameDeal, freeNow []models.GameDeal, err error) {
	url := "https://store.steampowered.com/api/featuredcategories/?l=brazilian&cc=br"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao criar requisição steam: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao chamar api steam: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("steam retornou status %d", resp.StatusCode)
	}

	var steamCat steamFeaturedCategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&steamCat); err != nil {
		return nil, nil, fmt.Errorf("erro ao decodificar resposta da steam: %w", err)
	}

	for _, item := range steamCat.Specials.Items {
		origPrice := float64(item.OriginalPrice) / 100.0
		finalPrice := float64(item.FinalPrice) / 100.0

		imageURL := item.LargeCapsuleImage
		if imageURL == "" {
			imageURL = item.HeaderImage
		}

		storeURL := fmt.Sprintf("https://store.steampowered.com/app/%d", item.ID)

		var promoEnd *time.Time
		if item.DiscountExpiration > 0 {
			exp := time.Unix(item.DiscountExpiration, 0)
			promoEnd = &exp
		}

		if item.DiscountPercent == 100 || (finalPrice == 0 && origPrice > 0) {
			freeNow = append(freeNow, models.GameDeal{
				ID:              fmt.Sprintf("steam_%d", item.ID),
				Title:           item.Name,
				Store:           "steam",
				Status:          "free_now",
				OriginalPrice:   origPrice,
				FinalPrice:      0,
				DiscountPercent: 100,
				Currency:        "BRL",
				ImageUrl:        imageURL,
				StoreUrl:        storeURL,
				EndDate:         promoEnd,
			})
		} else if item.Discounted && item.DiscountPercent > 0 {
			specials = append(specials, models.GameDeal{
				ID:              fmt.Sprintf("steam_%d", item.ID),
				Title:           item.Name,
				Store:           "steam",
				Status:          "deal",
				OriginalPrice:   origPrice,
				FinalPrice:      finalPrice,
				DiscountPercent: item.DiscountPercent,
				Currency:        "BRL",
				ImageUrl:        imageURL,
				StoreUrl:        storeURL,
				EndDate:         promoEnd,
			})
		}
	}

	return specials, freeNow, nil
}
