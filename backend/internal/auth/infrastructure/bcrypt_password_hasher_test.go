package infrastructure

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/auth/domain"
)

func TestBcryptPasswordHasher(t *testing.T) {
	ctx := context.Background()
	hasher := NewBcryptPasswordHasher(4)

	hash, err := hasher.Hash(ctx, "correct-password")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "" || hash == "correct-password" {
		t.Fatalf("expected bcrypt hash, got %q", hash)
	}

	if err := hasher.Compare(ctx, "correct-password", hash); err != nil {
		t.Fatalf("compare correct password: %v", err)
	}

	err = hasher.Compare(ctx, "wrong-password", hash)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestBcryptPasswordHasherRejectsEmptyPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher(4)

	_, err := hasher.Hash(context.Background(), " ")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
