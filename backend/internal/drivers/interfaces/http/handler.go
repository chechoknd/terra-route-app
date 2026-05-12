package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type CreateDriverUseCase interface {
	Execute(ctx context.Context, req dto.CreateDriverRequest) (*dto.Driver, error)
}

type UpdateDriverUseCase interface {
	Execute(ctx context.Context, req dto.UpdateDriverRequest) (*dto.Driver, error)
}

type GetDriverUseCase interface {
	Execute(ctx context.Context, req dto.GetDriverRequest) (*dto.Driver, error)
}

type ListDriversUseCase interface {
	Execute(ctx context.Context, req dto.ListDriversRequest) ([]dto.Driver, error)
}

type DeactivateDriverUseCase interface {
	Execute(ctx context.Context, req dto.DeactivateDriverRequest) error
}

type Handler struct {
	create     CreateDriverUseCase
	update     UpdateDriverUseCase
	get        GetDriverUseCase
	list       ListDriversUseCase
	deactivate DeactivateDriverUseCase
}

func NewHandler(
	create CreateDriverUseCase,
	update UpdateDriverUseCase,
	get GetDriverUseCase,
	list ListDriversUseCase,
	deactivate DeactivateDriverUseCase,
) *Handler {
	return &Handler{
		create:     create,
		update:     update,
		get:        get,
		list:       list,
		deactivate: deactivate,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	var req driverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	driver, err := h.create.Execute(r.Context(), dto.CreateDriverRequest{
		CompanyID:      claims.CompanyID,
		UserID:         req.UserID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocumentNumber: req.DocumentNumber,
		Phone:          req.Phone,
		Email:          req.Email,
		LicenseNumber:  req.LicenseNumber,
		Status:         req.Status,
	})
	if err != nil {
		writeDriverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, driverEnvelope{Driver: driverToResponse(*driver)})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	drivers, err := h.list.Execute(r.Context(), dto.ListDriversRequest{CompanyID: claims.CompanyID})
	if err != nil {
		writeDriverError(w, err)
		return
	}

	res := make([]driverResponse, 0, len(drivers))
	for _, driver := range drivers {
		res = append(res, driverToResponse(driver))
	}
	writeJSON(w, http.StatusOK, driversEnvelope{Drivers: res})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	driver, err := h.get.Execute(r.Context(), dto.GetDriverRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeDriverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, driverEnvelope{Driver: driverToResponse(*driver)})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	var req driverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	driver, err := h.update.Execute(r.Context(), dto.UpdateDriverRequest{
		CompanyID:      claims.CompanyID,
		ID:             r.PathValue("id"),
		UserID:         req.UserID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocumentNumber: req.DocumentNumber,
		Phone:          req.Phone,
		Email:          req.Email,
		LicenseNumber:  req.LicenseNumber,
		Status:         req.Status,
	})
	if err != nil {
		writeDriverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, driverEnvelope{Driver: driverToResponse(*driver)})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	err := h.deactivate.Execute(r.Context(), dto.DeactivateDriverRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeDriverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func driverToResponse(driver dto.Driver) driverResponse {
	return driverResponse{
		ID:             driver.ID,
		CompanyID:      driver.CompanyID,
		UserID:         driver.UserID,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		DocumentNumber: driver.DocumentNumber,
		Phone:          driver.Phone,
		Email:          driver.Email,
		LicenseNumber:  driver.LicenseNumber,
		Status:         driver.Status,
		CreatedAt:      driver.CreatedAt,
		UpdatedAt:      driver.UpdatedAt,
	}
}

func writeDriverError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidDriver):
		writeError(w, http.StatusBadRequest, "invalid_driver")
	case errors.Is(err, domain.ErrDriverNotFound):
		writeError(w, http.StatusNotFound, "driver_not_found")
	case errors.Is(err, domain.ErrDriverAlreadyExists):
		writeError(w, http.StatusConflict, "driver_already_exists")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
