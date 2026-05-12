package http

import (
	"errors"
	"net/http"
	"strings"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
)

func AuthMiddleware(tokens authdomain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				writeError(w, http.StatusUnauthorized, "missing_or_malformed_token")
				return
			}

			claims, err := tokens.Validate(r.Context(), tokenValue)
			if errors.Is(err, authdomain.ErrInvalidToken) {
				writeError(w, http.StatusUnauthorized, "invalid_token")
				return
			}
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal_error")
				return
			}

			next.ServeHTTP(w, r.WithContext(contextWithAuthClaims(r.Context(), claims)))
		})
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}
	return parts[1], true
}
