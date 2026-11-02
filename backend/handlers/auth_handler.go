package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ppdm/db"
	"ppdm/middlewares"
	"ppdm/models"
	"ppdm/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	UserRepo *db.UserRepository
}

func NewAuthHandler(userRepo *db.UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo: userRepo,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 6 {
		http.Error(w, `{"error": "Email é obrigatório e senha deve ter no mínimo 6 caracteres"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// verificar se ja existe ne
	existing, err := h.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, `{"error": "Erro ao verificar usuário"}`, http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, `{"error": "Email já cadastrado"}`, http.StatusConflict)
		return
	}

	// cripto aqui
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error": "Erro ao processar senha"}`, http.StatusInternalServerError)
		return
	}

	newUser := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	insertedID, err := h.UserRepo.Create(ctx, newUser)
	if err != nil {
		http.Error(w, `{"error": "Erro ao salvar usuário no banco"}`, http.StatusInternalServerError)
		return
	}
	newUser.ID = insertedID

	token, err := utils.GenerateToken(insertedID.Hex(), newUser.Email, newUser.Role)
	if err != nil {
		http.Error(w, `{"error": "Erro ao gerar token"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  newUser.ToResponse(),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	ctx := r.Context()

	user, err := h.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		http.Error(w, `{"error": "Email ou senha inválidos"}`, http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, `{"error": "Email ou senha inválidos"}`, http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error": "Erro ao gerar token de autenticação"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := middlewares.GetUserClaims(r.Context())
	if !ok || claims == nil {
		http.Error(w, `{"error": "Não autenticado"}`, http.StatusUnauthorized)
		return
	}

	objID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, `{"error": "ID de usuário inválido"}`, http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.FindByID(r.Context(), objID)
	if err != nil || user == nil {
		http.Error(w, `{"error": "Usuário não encontrado"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user.ToResponse())
}

