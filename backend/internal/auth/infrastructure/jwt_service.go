package infrastructure

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/terraroute/terra-route/backend/internal/auth/domain"
)

type JWTService struct {
	secret     []byte
	expiration time.Duration
	now        func() time.Time
}

type jwtClaims struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiration string) (*JWTService, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, domain.ErrInvalidToken
	}

	duration, err := time.ParseDuration(expiration)
	if err != nil || duration <= 0 {
		return nil, domain.ErrInvalidToken
	}

	return &JWTService{
		secret:     []byte(secret),
		expiration: duration,
		now:        time.Now,
	}, nil
}

func (s *JWTService) Generate(ctx context.Context, subject domain.TokenSubject) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(subject.UserID) == "" ||
		strings.TrimSpace(subject.CompanyID) == "" ||
		strings.TrimSpace(subject.Role) == "" {
		return "", domain.ErrInvalidToken
	}

	now := s.now().UTC()
	expiresAt := now.Add(s.expiration)
	claims := jwtClaims{
		UserID:    subject.UserID,
		CompanyID: subject.CompanyID,
		Role:      subject.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (s *JWTService) Validate(ctx context.Context, tokenValue string) (*domain.TokenClaims, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tokenValue) == "" {
		return nil, domain.ErrInvalidToken
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, domain.ErrInvalidToken
		}
		return s.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithLeeway(5*time.Second), jwt.WithTimeFunc(s.now))
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	if token == nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	if claims.Subject != claims.UserID {
		return nil, domain.ErrInvalidToken
	}
	if strings.TrimSpace(claims.UserID) == "" ||
		strings.TrimSpace(claims.CompanyID) == "" ||
		strings.TrimSpace(claims.Role) == "" ||
		claims.ExpiresAt == nil {
		return nil, domain.ErrInvalidToken
	}
	if claims.ExpiresAt.Time.Before(s.now()) {
		return nil, domain.ErrInvalidToken
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		CompanyID: claims.CompanyID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
