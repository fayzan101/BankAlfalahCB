package middleware

import (
	"net/http"
	"strings"

	"bank-ai-chatbot/internal/api/handlers"
	"bank-ai-chatbot/internal/security"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
)

type AuthMiddleware struct {
	tokens *security.TokenManager
}

func NewAuthMiddleware(tokens *security.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			response.Error(w, err)
			return
		}

		claims, err := m.tokens.Validate(token)
		if err != nil {
			response.Error(w, apperrors.Unauthorized("invalid or expired token"))
			return
		}

		ctx := handlers.WithUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", apperrors.Unauthorized("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", apperrors.Unauthorized("invalid authorization header format")
	}

	return strings.TrimSpace(parts[1]), nil
}
