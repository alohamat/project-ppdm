package routes

import (
	"github.com/gorilla/mux"
	"ppdm/handlers"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterHandler)

	return router
}