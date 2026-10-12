package handlers

import (
	"net/http"
	"encoding/json"
	"ppdm/models"
	"log"
)

func LoginHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		http.Error(w, "JSON inválido ou corpo vazio", http.StatusBadRequest)
		return
	}
	log.Println(req.Email)
	log.Println(req.Password)
}