package infrastructure

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/terraroute/terra-route/backend/internal/auth/domain"
)

const DefaultBcryptCost = bcrypt.DefaultCost

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost == 0 {
		cost = DefaultBcryptCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

func (h *BcryptPasswordHasher) Hash(ctx context.Context, password string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(password) == "" {
		return "", domain.ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *BcryptPasswordHasher) Compare(ctx context.Context, password string, passwordHash string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(password) == "" || strings.TrimSpace(passwordHash) == "" {
		return domain.ErrInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return domain.ErrInvalidCredentials
	}
	if err != nil {
		return err
	}

	return nil
}
