package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/terraroute/terra-route/backend/internal/auth/domain"
)

func TestJWTServiceGenerateAndValidate(t *testing.T) {
	service, err := NewJWTService("test-secret", "15m")
	if err != nil {
		t.Fatalf("new jwt service: %v", err)
	}
	fixedNow := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return fixedNow }

	token, err := service.Generate(context.Background(), domain.TokenSubject{
		UserID:    "user-1",
		CompanyID: "company-1",
		Role:      "operator",
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}

	claims, err := service.Validate(context.Background(), token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}

	if claims.UserID != "user-1" {
		t.Fatalf("expected user_id user-1, got %q", claims.UserID)
	}
	if claims.CompanyID != "company-1" {
		t.Fatalf("expected company_id company-1, got %q", claims.CompanyID)
	}
	if claims.Role != "operator" {
		t.Fatalf("expected role operator, got %q", claims.Role)
	}
	if !claims.ExpiresAt.Equal(fixedNow.Add(15 * time.Minute)) {
		t.Fatalf("unexpected expiration: %s", claims.ExpiresAt)
	}
}

func TestJWTServiceRejectsInvalidSecretOrExpiration(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		expiration string
	}{
		{name: "empty secret", secret: "", expiration: "15m"},
		{name: "empty expiration", secret: "test-secret", expiration: ""},
		{name: "invalid expiration", secret: "test-secret", expiration: "soon"},
		{name: "zero expiration", secret: "test-secret", expiration: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewJWTService(tt.secret, tt.expiration)
			if !errors.Is(err, domain.ErrInvalidToken) {
				t.Fatalf("expected ErrInvalidToken, got %v", err)
			}
		})
	}
}

func TestJWTServiceRejectsInvalidSubject(t *testing.T) {
	service, err := NewJWTService("test-secret", "15m")
	if err != nil {
		t.Fatalf("new jwt service: %v", err)
	}

	_, err = service.Generate(context.Background(), domain.TokenSubject{
		UserID: "user-1",
		Role:   "operator",
	})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTServiceRejectsExpiredToken(t *testing.T) {
	service, err := NewJWTService("test-secret", "15m")
	if err != nil {
		t.Fatalf("new jwt service: %v", err)
	}
	fixedNow := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return fixedNow }

	token, err := service.Generate(context.Background(), domain.TokenSubject{
		UserID:    "user-1",
		CompanyID: "company-1",
		Role:      "operator",
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	service.now = func() time.Time { return fixedNow.Add(16 * time.Minute) }
	_, err = service.Validate(context.Background(), token)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTServiceRejectsWrongSecret(t *testing.T) {
	service, err := NewJWTService("test-secret", "15m")
	if err != nil {
		t.Fatalf("new jwt service: %v", err)
	}

	token, err := service.Generate(context.Background(), domain.TokenSubject{
		UserID:    "user-1",
		CompanyID: "company-1",
		Role:      "operator",
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	otherService, err := NewJWTService("other-secret", "15m")
	if err != nil {
		t.Fatalf("new other jwt service: %v", err)
	}

	_, err = otherService.Validate(context.Background(), token)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTServiceRejectsUnexpectedSigningMethod(t *testing.T) {
	service, err := NewJWTService("test-secret", "15m")
	if err != nil {
		t.Fatalf("new jwt service: %v", err)
	}

	claims := jwt.MapClaims{
		"user_id":    "user-1",
		"company_id": "company-1",
		"role":       "operator",
		"sub":        "user-1",
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenValue, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none token: %v", err)
	}

	_, err = service.Validate(context.Background(), tokenValue)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
