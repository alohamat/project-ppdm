package services

import (
	"context"
	"testing"
	"time"
)

func TestDealsService(t *testing.T) {
	service := NewDealsService()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	deals, err := service.GetDeals(ctx)
	if err != nil {
		t.Fatalf("Erro ao buscar ofertas: %v", err)
	}

	if deals == nil {
		t.Fatal("Deals response é nil")
	}

	t.Logf("Jogos grátis agora: %d", len(deals.FreeNow))
	for _, g := range deals.FreeNow {
		t.Logf("- [GRÁTIS AGORA] (%s) %s -> %s", g.Store, g.Title, g.StoreUrl)
	}

	t.Logf("Jogos grátis futuros: %d", len(deals.FreeUpcoming))
	for _, g := range deals.FreeUpcoming {
		t.Logf("- [GRÁTIS EM BREVE] (%s) %s -> %s", g.Store, g.Title, g.StoreUrl)
	}

	t.Logf("Super ofertas: %d", len(deals.BigDeals))
	for i, g := range deals.BigDeals {
		if i >= 5 {
			break
		}
		t.Logf("- [OFERTA -%d%%] (%s) %s | De R$%.2f por R$%.2f -> %s", g.DiscountPercent, g.Store, g.Title, g.OriginalPrice, g.FinalPrice, g.StoreUrl)
	}

	if len(deals.FreeNow) == 0 && len(deals.FreeUpcoming) == 0 && len(deals.BigDeals) == 0 {
		t.Error("Nenhum jogo encontrado")
	}
}
