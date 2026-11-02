package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	objID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID de usuário inválido"}`, http.StatusBadRequest)
		return
	}

	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	req.Role = strings.TrimSpace(strings.ToLower(req.Role))
	if req.Role != "admin" && req.Role != "user" {
		http.Error(w, `{"error": "Cargo inválido. Escolha 'admin' ou 'user'"}`, http.StatusBadRequest)
		return
	}

	targetUser, err := h.UserRepo.FindByID(r.Context(), objID)
	if err != nil {
		http.Error(w, `{"error": "Erro ao verificar usuário"}`, http.StatusInternalServerError)
		return
	}
	if targetUser == nil {
		http.Error(w, `{"error": "Usuário não encontrado"}`, http.StatusNotFound)
		return
	}

	if err := h.UserRepo.UpdateRole(r.Context(), objID, req.Role); err != nil {
		http.Error(w, `{"error": "Erro ao atualizar cargo do usuário"}`, http.StatusInternalServerError)
		return
	}

	targetUser.Role = req.Role
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Cargo atualizado com sucesso",
		"user":    targetUser.ToResponse(),
	})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	objID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID de usuário inválido"}`, http.StatusBadRequest)
		return
	}

	targetUser, err := h.UserRepo.FindByID(r.Context(), objID)
	if err != nil {
		http.Error(w, `{"error": "Erro ao verificar usuário"}`, http.StatusInternalServerError)
		return
	}
	if targetUser == nil {
		http.Error(w, `{"error": "Usuário não encontrado"}`, http.StatusNotFound)
		return
	}

	if err := h.UserRepo.Delete(r.Context(), objID); err != nil {
		http.Error(w, `{"error": "Erro ao deletar usuário"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuário deletado com sucesso",
	})
}

