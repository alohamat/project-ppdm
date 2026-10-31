package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"ppdm/db"
	"ppdm/handlers"
	"ppdm/middlewares"
	"ppdm/services"
)

func InitRouter(userRepo *db.UserRepository, dealsService *services.DealsService) *mux.Router {
	router := mux.NewRouter()

	authHandler := handlers.NewAuthHandler(userRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)
	dealsHandler := handlers.NewDealsHandler(dealsService)

	// Rotas Públicas - Autenticação
	router.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	// Rotas Públicas - Jogos Grátis e Ofertas (Epic Games e Steam)
	router.HandleFunc("/api/deals", dealsHandler.GetAllDeals).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/deals/free", dealsHandler.GetFreeDeals).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/deals/specials", dealsHandler.GetSpecialDeals).Methods(http.MethodGet, http.MethodOptions)

	// Subrouter Protegido (Autenticado via AuthMiddleware)
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	protected.HandleFunc("/me", authHandler.Me).Methods(http.MethodGet, http.MethodOptions)

	// Subrouter de Admin (Autenticado + AdminMiddleware)
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middlewares.AdminMiddleware)

	admin.HandleFunc("/users", adminHandler.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	admin.HandleFunc("/stats", adminHandler.GetStats).Methods(http.MethodGet, http.MethodOptions)

	// Middleware de CORS
	router.Use(middlewares.CorsMiddleware)

	return router
}
