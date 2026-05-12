package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/terraroute/terra-route/backend/internal/auth/application/dto"
	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, companyID string, email string) (*userdomain.User, error)
}

type LoginUseCase struct {
	users          UserRepository
	passwordHasher authdomain.PasswordHasher
	tokenService   authdomain.TokenService
}

func NewLoginUseCase(users UserRepository, passwordHasher authdomain.PasswordHasher, tokenService authdomain.TokenService) *LoginUseCase {
	return &LoginUseCase{
		users:          users,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	companyID := strings.TrimSpace(req.CompanyID)
	email := strings.TrimSpace(req.Email)
	password := req.Password

	if companyID == "" || email == "" || strings.TrimSpace(password) == "" {
		return nil, authdomain.ErrInvalidCredentials
	}

	user, err := uc.users.GetByEmail(ctx, companyID, email)
	if errors.Is(err, userdomain.ErrUserNotFound) || errors.Is(err, userdomain.ErrInvalidUser) {
		return nil, authdomain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if user.Status != userdomain.UserStatusActive {
		return nil, authdomain.ErrInactiveUser
	}

	if err := uc.passwordHasher.Compare(ctx, password, user.PasswordHash); err != nil {
		if errors.Is(err, authdomain.ErrInvalidCredentials) {
			return nil, authdomain.ErrInvalidCredentials
		}
		return nil, err
	}

	accessToken, err := uc.tokenService.Generate(ctx, authdomain.TokenSubject{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Role:      user.Role,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User: dto.AuthenticatedUser{
			ID:        user.ID,
			CompanyID: user.CompanyID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}, nil
}
