package middleware

import (
	"context"
	"net/http"
	
	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/services/auth"
)

type AuthMiddleware struct {
	authService *auth.AuthService
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		authService: auth.NewAuthService(cfg),
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter token do parâmetro de consulta
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		// Validar token
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Armazenar informações do usuário no contexto
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}