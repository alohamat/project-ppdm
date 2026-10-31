package handlers

import (
	"encoding/json"
	"net/http"

	"ppdm/services"
)

type DealsHandler struct {
	dealsService *services.DealsService
}

func NewDealsHandler(dealsService *services.DealsService) *DealsHandler {
	return &DealsHandler{
		dealsService: dealsService,
	}
}

func (h *DealsHandler) GetAllDeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deals, err := h.dealsService.GetDeals(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Falha ao obter ofertas de jogos"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deals)
}

func (h *DealsHandler) GetFreeDeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deals, err := h.dealsService.GetDeals(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Falha ao obter jogos gratuitos"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"free_now":      deals.FreeNow,
		"free_upcoming": deals.FreeUpcoming,
		"updated_at":    deals.UpdatedAt,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *DealsHandler) GetSpecialDeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deals, err := h.dealsService.GetDeals(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Falha ao obter promoções"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"big_deals":  deals.BigDeals,
		"updated_at": deals.UpdatedAt,
	}

	json.NewEncoder(w).Encode(response)
}
