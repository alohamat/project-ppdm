package handlers

import (
	"encoding/json"
	"net/http"

	"ppdm/db"
	"ppdm/models"
)

type AdminHandler struct {
	UserRepo *db.UserRepository
}

func NewAdminHandler(userRepo *db.UserRepository) *AdminHandler {
	return &AdminHandler{
		UserRepo: userRepo,
	}
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Falha ao buscar usuários"}`, http.StatusInternalServerError)
		return
	}

	responseList := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		responseList = append(responseList, u.ToResponse())
	}

	json.NewEncoder(w).Encode(responseList)
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Falha ao obter estatísticas"}`, http.StatusInternalServerError)
		return
	}

	totalUsers := len(users)
	totalAdmins := 0
	for _, u := range users {
		if u.Role == "admin" {
			totalAdmins++
		}
	}

	stats := map[string]interface{}{
		"total_users":  totalUsers,
		"total_admins": totalAdmins,
		"total_normal": totalUsers - totalAdmins,
	}

	json.NewEncoder(w).Encode(stats)
}

