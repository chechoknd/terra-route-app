package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/auth/application/dto"
	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestHandlerLogin(t *testing.T) {
	handler := NewHandler(
		&fakeLoginUseCase{
			response: &dto.LoginResponse{
				AccessToken: "token",
				TokenType:   "Bearer",
				User: dto.AuthenticatedUser{
					ID:        "user-1",
					CompanyID: "company-1",
					Email:     "operator@example.com",
					FullName:  "Operator",
					Role:      userdomain.UserRoleOperator,
					Status:    userdomain.UserStatusActive,
				},
			},
		},
		&fakeTokenService{},
		&fakeUserReader{},
	)

	body := bytes.NewBufferString(`{"company_id":"company-1","email":"operator@example.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res loginResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.AccessToken != "token" || res.TokenType != "Bearer" {
		t.Fatalf("unexpected token response: %+v", res)
	}
	if res.User.ID != "user-1" {
		t.Fatalf("expected user-1, got %q", res.User.ID)
	}
}

func TestHandlerLoginRejectsInvalidCredentials(t *testing.T) {
	handler := NewHandler(&fakeLoginUseCase{err: authdomain.ErrInvalidCredentials}, &fakeTokenService{}, &fakeUserReader{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{"company_id":"company-1"}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestHandlerMe(t *testing.T) {
	handler := NewHandler(
		&fakeLoginUseCase{},
		&fakeTokenService{
			claims: &authdomain.TokenClaims{
				UserID:    "user-1",
				CompanyID: "company-1",
				Role:      userdomain.UserRoleOperator,
			},
		},
		&fakeUserReader{
			user: &userdomain.User{
				ID:        "user-1",
				CompanyID: "company-1",
				Email:     "operator@example.com",
				FullName:  "Operator",
				Role:      userdomain.UserRoleOperator,
				Status:    userdomain.UserStatusActive,
			},
		},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res meResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.User.ID != "user-1" {
		t.Fatalf("expected user-1, got %q", res.User.ID)
	}
}

func TestHandlerMeRejectsMissingToken(t *testing.T) {
	handler := NewHandler(&fakeLoginUseCase{}, &fakeTokenService{}, &fakeUserReader{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestHandlerMeRejectsInvalidToken(t *testing.T) {
	handler := NewHandler(&fakeLoginUseCase{}, &fakeTokenService{err: authdomain.ErrInvalidToken}, &fakeUserReader{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

type fakeLoginUseCase struct {
	response *dto.LoginResponse
	err      error
}

func (uc *fakeLoginUseCase) Execute(_ context.Context, _ dto.LoginRequest) (*dto.LoginResponse, error) {
	return uc.response, uc.err
}

type fakeTokenService struct {
	claims *authdomain.TokenClaims
	err    error
}

func (s *fakeTokenService) Generate(_ context.Context, _ authdomain.TokenSubject) (string, error) {
	return "token", nil
}

func (s *fakeTokenService) Validate(_ context.Context, _ string) (*authdomain.TokenClaims, error) {
	return s.claims, s.err
}

type fakeUserReader struct {
	user *userdomain.User
	err  error
}

func (r *fakeUserReader) GetByID(_ context.Context, _ string, _ string) (*userdomain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.user == nil {
		return nil, errors.New("missing fake user")
	}
	return r.user, nil
}
