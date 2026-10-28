package middlewares

import (
	"context"
	"net/http"
	"strings"

	"ppdm/utils"
)

type contextKey string

const (
	UserClaimsKey contextKey = "userClaims"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Token de autorização não fornecido"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error": "Formato de autorização inválido. Use 'Bearer <token>'"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": "Token inválido ou expirado"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
		if !ok || claims == nil {
			http.Error(w, `{"error": "Usuário não autenticado"}`, http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, `{"error": "Acesso negado: privilégios de administrador necessários"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserClaims(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*utils.Claims)
	return claims, ok
}

