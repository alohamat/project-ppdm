package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"ppdm/db"
	"ppdm/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, utilizando variáveis de ambiente do sistema.")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI não encontrada no arquivo .env ou nas variáveis de ambiente")
	}

	database, err := db.ConnectDB(mongoURI, "ppdm")
	if err != nil {
		log.Fatalf("Erro de conexão com MongoDB: %v", err)
	}
	defer database.Close()
	log.Println("MongoDB conectado com sucesso!")

	userRepo := db.NewUserRepository(database.Collection("users"))
	router := routes.InitRouter(userRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor backend rodando em http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
