package main

import (
	"ppdm/db"
	"ppdm/routes"
	"github.com/joho/godotenv"

	"net/http"
	"os"
	"log"
)

func main() {
	godotenv.Load()
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatal("uri not found in .env")
	}

	db, err := db.ConnectDB(MONGO_URI, "ppdm")
	if err != nil {
		log.Fatalf("Erro de conexão: %v", err)
	}
	defer db.Close()

	router := routes.InitRouter()

	log.Println("Rodando em :8080")
	log.Fatalln(http.ListenAndServe(":8080", router))
}