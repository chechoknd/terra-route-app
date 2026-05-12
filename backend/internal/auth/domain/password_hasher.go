package domain

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, password string, passwordHash string) error
}
