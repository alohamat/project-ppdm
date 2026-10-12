package handlers

import (
	"net/http"
	"encoding/json"
	"ppdm/models"
	"log"
)

func RegisterHandler (w http.ResponseWriter, r *http.Request) {
	var Req models.User
	err := json.NewDecoder(r.Body).Decode(&Req)
	if (err != nil) {
		log.Fatalln("Erro decodificando JSON")
	}
	log.Println(Req.Email)
	log.Println(Req.Password)
}