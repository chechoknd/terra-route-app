package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/terraroute/terra-route/backend/internal/auth/application/dto"
	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

type LoginUseCase interface {
	Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type UserReader interface {
	GetByID(ctx context.Context, companyID string, id string) (*userdomain.User, error)
}

type Handler struct {
	login  LoginUseCase
	tokens authdomain.TokenService
	users  UserReader
}

func NewHandler(login LoginUseCase, tokens authdomain.TokenService, users UserReader) *Handler {
	return &Handler{
		login:  login,
		tokens: tokens,
		users:  users,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	res, err := h.login.Execute(r.Context(), dto.LoginRequest{
		CompanyID: req.CompanyID,
		Email:     req.Email,
		Password:  req.Password,
	})
	if errors.Is(err, authdomain.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}
	if errors.Is(err, authdomain.ErrInactiveUser) {
		writeError(w, http.StatusForbidden, "inactive_user")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		AccessToken: res.AccessToken,
		TokenType:   res.TokenType,
		User: userResponse{
			ID:        res.User.ID,
			CompanyID: res.User.CompanyID,
			Email:     res.User.Email,
			FullName:  res.User.FullName,
			Role:      res.User.Role,
			Status:    res.User.Status,
			CreatedAt: res.User.CreatedAt,
			UpdatedAt: res.User.UpdatedAt,
		},
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid_token")
		return
	}

	user, err := h.users.GetByID(r.Context(), claims.CompanyID, claims.UserID)
	if errors.Is(err, userdomain.ErrUserNotFound) || errors.Is(err, userdomain.ErrInvalidUser) {
		writeError(w, http.StatusUnauthorized, "invalid_token")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}
	if user.Status != userdomain.UserStatusActive {
		writeError(w, http.StatusForbidden, "inactive_user")
		return
	}

	writeJSON(w, http.StatusOK, meResponse{
		User: userResponse{
			ID:        user.ID,
			CompanyID: user.CompanyID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
