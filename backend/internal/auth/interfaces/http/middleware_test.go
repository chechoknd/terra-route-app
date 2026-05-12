package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
)

func TestAuthMiddlewareStoresClaimsInContext(t *testing.T) {
	tokens := &fakeTokenService{
		claims: &authdomain.TokenClaims{
			UserID:    "user-1",
			CompanyID: "company-1",
			Role:      "operator",
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := AuthClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("expected auth claims in context")
		}
		if claims.UserID != "user-1" || claims.CompanyID != "company-1" || claims.Role != "operator" {
			t.Fatalf("unexpected claims: %+v", claims)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	AuthMiddleware(tokens)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if tokens.validatedToken != "valid-token" {
		t.Fatalf("expected token to be validated without logging/exposing it")
	}
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	AuthMiddleware(&fakeTokenService{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "missing_or_malformed_token")
}

func TestAuthMiddlewareRejectsMalformedToken(t *testing.T) {
	tests := []string{
		"token",
		"Basic token",
		"Bearer",
		"Bearer ",
	}

	for _, header := range tests {
		t.Run(header, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", header)

			AuthMiddleware(&fakeTokenService{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("next handler should not run")
			})).ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected status 401, got %d", rec.Code)
			}
			assertJSONError(t, rec.Body.String(), "missing_or_malformed_token")
		})
	}
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	AuthMiddleware(&fakeTokenService{err: authdomain.ErrInvalidToken})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "invalid_token")
}

func TestAuthMiddlewareReturnsInternalErrorForUnexpectedTokenServiceError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer token")

	AuthMiddleware(&fakeTokenService{err: errors.New("token service unavailable")})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "internal_error")
}

func assertJSONError(t *testing.T, body string, want string) {
	t.Helper()
	if body != `{"error":"`+want+`"}`+"\n" {
		t.Fatalf("expected JSON error %q, got %q", want, body)
	}
}
