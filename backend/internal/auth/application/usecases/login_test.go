package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/auth/application/dto"
	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestLoginUseCaseAuthenticatesActiveUser(t *testing.T) {
	users := &fakeUserRepository{
		user: &userdomain.User{
			ID:           "user-1",
			CompanyID:    "company-1",
			Email:        "operator@example.com",
			FullName:     "Ops User",
			Role:         userdomain.UserRoleOperator,
			Status:       userdomain.UserStatusActive,
			PasswordHash: "valid-hash",
		},
	}
	hasher := &fakePasswordHasher{}
	tokens := &fakeTokenService{token: "access-token"}
	uc := NewLoginUseCase(users, hasher, tokens)

	res, err := uc.Execute(context.Background(), dto.LoginRequest{
		CompanyID: "company-1",
		Email:     "operator@example.com",
		Password:  "secret",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if res.User.ID != "user-1" {
		t.Fatalf("expected authenticated user id user-1, got %q", res.User.ID)
	}
	if res.User.Email != "operator@example.com" {
		t.Fatalf("expected email operator@example.com, got %q", res.User.Email)
	}
	if hasher.password != "secret" || hasher.passwordHash != "valid-hash" {
		t.Fatalf("password hasher was called with unexpected values")
	}
	if res.AccessToken != "access-token" {
		t.Fatalf("expected access token, got %q", res.AccessToken)
	}
	if tokens.subject.UserID != "user-1" || tokens.subject.CompanyID != "company-1" || tokens.subject.Role != userdomain.UserRoleOperator {
		t.Fatalf("token service was called with unexpected subject")
	}
}

func TestLoginUseCaseRejectsUnknownUser(t *testing.T) {
	uc := NewLoginUseCase(&fakeUserRepository{err: userdomain.ErrUserNotFound}, &fakePasswordHasher{}, &fakeTokenService{})

	_, err := uc.Execute(context.Background(), dto.LoginRequest{
		CompanyID: "company-1",
		Email:     "missing@example.com",
		Password:  "secret",
	})
	if !errors.Is(err, authdomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUseCaseRejectsWrongPassword(t *testing.T) {
	uc := NewLoginUseCase(
		&fakeUserRepository{
			user: &userdomain.User{
				ID:           "user-1",
				CompanyID:    "company-1",
				Email:        "operator@example.com",
				FullName:     "Ops User",
				Role:         userdomain.UserRoleOperator,
				Status:       userdomain.UserStatusActive,
				PasswordHash: "valid-hash",
			},
		},
		&fakePasswordHasher{err: authdomain.ErrInvalidCredentials},
		&fakeTokenService{},
	)

	_, err := uc.Execute(context.Background(), dto.LoginRequest{
		CompanyID: "company-1",
		Email:     "operator@example.com",
		Password:  "wrong",
	})
	if !errors.Is(err, authdomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUseCaseRejectsInactiveUser(t *testing.T) {
	uc := NewLoginUseCase(
		&fakeUserRepository{
			user: &userdomain.User{
				ID:           "user-1",
				CompanyID:    "company-1",
				Email:        "operator@example.com",
				FullName:     "Ops User",
				Role:         userdomain.UserRoleOperator,
				Status:       userdomain.UserStatusInactive,
				PasswordHash: "valid-hash",
			},
		},
		&fakePasswordHasher{},
		&fakeTokenService{},
	)

	_, err := uc.Execute(context.Background(), dto.LoginRequest{
		CompanyID: "company-1",
		Email:     "operator@example.com",
		Password:  "secret",
	})
	if !errors.Is(err, authdomain.ErrInactiveUser) {
		t.Fatalf("expected ErrInactiveUser, got %v", err)
	}
}

func TestLoginUseCaseRequiresCompanyScope(t *testing.T) {
	uc := NewLoginUseCase(&fakeUserRepository{}, &fakePasswordHasher{}, &fakeTokenService{})

	_, err := uc.Execute(context.Background(), dto.LoginRequest{
		Email:    "operator@example.com",
		Password: "secret",
	})
	if !errors.Is(err, authdomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

type fakeUserRepository struct {
	user *userdomain.User
	err  error
}

func (r *fakeUserRepository) GetByEmail(_ context.Context, _ string, _ string) (*userdomain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.user, nil
}

type fakePasswordHasher struct {
	password     string
	passwordHash string
	err          error
}

func (h *fakePasswordHasher) Hash(_ context.Context, password string) (string, error) {
	return "hashed-" + password, nil
}

func (h *fakePasswordHasher) Compare(_ context.Context, password string, passwordHash string) error {
	h.password = password
	h.passwordHash = passwordHash
	return h.err
}

type fakeTokenService struct {
	token   string
	subject authdomain.TokenSubject
	err     error
}

func (s *fakeTokenService) Generate(_ context.Context, subject authdomain.TokenSubject) (string, error) {
	s.subject = subject
	return s.token, s.err
}

func (s *fakeTokenService) Validate(_ context.Context, _ string) (*authdomain.TokenClaims, error) {
	return nil, nil
}
