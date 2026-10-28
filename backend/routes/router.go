package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"ppdm/db"
	"ppdm/handlers"
	"ppdm/middlewares"
)

func InitRouter(userRepo *db.UserRepository) *mux.Router {
	router := mux.NewRouter()

	authHandler := handlers.NewAuthHandler(userRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)

	// rotas publicas
	router.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	protected.HandleFunc("/me", authHandler.Me).Methods(http.MethodGet, http.MethodOptions)

	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middlewares.AdminMiddleware)

	admin.HandleFunc("/users", adminHandler.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	admin.HandleFunc("/stats", adminHandler.GetStats).Methods(http.MethodGet, http.MethodOptions)

	router.Use(middlewares.CorsMiddleware)

	return router
}
