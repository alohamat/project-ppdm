package main

import (
	"net/http"
	"ppdm/routes"
	"log"
)

func main() {
	router := routes.InitRouter()

	log.Println("Rodando em :8080")
	log.Fatalln(http.ListenAndServe(":8080", router))
}