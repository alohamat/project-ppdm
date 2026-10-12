package routes

import (
	"github.com/gorilla/mux"
	"ppdm/handlers"
	"ppdm/middlewares"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterHandler)
	router.HandleFunc("/login", handlers.LoginHandler)

	router.Use(middlewares.CorsMiddleware)

	return router
}