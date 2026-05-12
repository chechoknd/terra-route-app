package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type CreateVehicleUseCase interface {
	Execute(ctx context.Context, req dto.CreateVehicleRequest) (*dto.Vehicle, error)
}

type UpdateVehicleUseCase interface {
	Execute(ctx context.Context, req dto.UpdateVehicleRequest) (*dto.Vehicle, error)
}

type GetVehicleUseCase interface {
	Execute(ctx context.Context, req dto.GetVehicleRequest) (*dto.Vehicle, error)
}

type ListVehiclesUseCase interface {
	Execute(ctx context.Context, req dto.ListVehiclesRequest) ([]dto.Vehicle, error)
}

type DeactivateVehicleUseCase interface {
	Execute(ctx context.Context, req dto.DeactivateVehicleRequest) error
}

type Handler struct {
	create     CreateVehicleUseCase
	update     UpdateVehicleUseCase
	get        GetVehicleUseCase
	list       ListVehiclesUseCase
	deactivate DeactivateVehicleUseCase
}

func NewHandler(
	create CreateVehicleUseCase,
	update UpdateVehicleUseCase,
	get GetVehicleUseCase,
	list ListVehiclesUseCase,
	deactivate DeactivateVehicleUseCase,
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

	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	vehicle, err := h.create.Execute(r.Context(), dto.CreateVehicleRequest{
		CompanyID:    claims.CompanyID,
		Plate:        req.Plate,
		InternalCode: req.InternalCode,
		VehicleType:  req.VehicleType,
		Brand:        req.Brand,
		Model:        req.Model,
		Capacity:     req.Capacity,
		Status:       req.Status,
	})
	if err != nil {
		writeVehicleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, vehicleEnvelope{Vehicle: vehicleToResponse(*vehicle)})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	vehicles, err := h.list.Execute(r.Context(), dto.ListVehiclesRequest{CompanyID: claims.CompanyID})
	if err != nil {
		writeVehicleError(w, err)
		return
	}

	res := make([]vehicleResponse, 0, len(vehicles))
	for _, vehicle := range vehicles {
		res = append(res, vehicleToResponse(vehicle))
	}
	writeJSON(w, http.StatusOK, vehiclesEnvelope{Vehicles: res})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	vehicle, err := h.get.Execute(r.Context(), dto.GetVehicleRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeVehicleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, vehicleEnvelope{Vehicle: vehicleToResponse(*vehicle)})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	vehicle, err := h.update.Execute(r.Context(), dto.UpdateVehicleRequest{
		CompanyID:    claims.CompanyID,
		ID:           r.PathValue("id"),
		Plate:        req.Plate,
		InternalCode: req.InternalCode,
		VehicleType:  req.VehicleType,
		Brand:        req.Brand,
		Model:        req.Model,
		Capacity:     req.Capacity,
		Status:       req.Status,
	})
	if err != nil {
		writeVehicleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, vehicleEnvelope{Vehicle: vehicleToResponse(*vehicle)})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	err := h.deactivate.Execute(r.Context(), dto.DeactivateVehicleRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeVehicleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func vehicleToResponse(vehicle dto.Vehicle) vehicleResponse {
	return vehicleResponse{
		ID:           vehicle.ID,
		CompanyID:    vehicle.CompanyID,
		Plate:        vehicle.Plate,
		InternalCode: vehicle.InternalCode,
		VehicleType:  vehicle.VehicleType,
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Capacity:     vehicle.Capacity,
		Status:       vehicle.Status,
		CreatedAt:    vehicle.CreatedAt,
		UpdatedAt:    vehicle.UpdatedAt,
	}
}

func writeVehicleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidVehicle):
		writeError(w, http.StatusBadRequest, "invalid_vehicle")
	case errors.Is(err, domain.ErrVehicleNotFound):
		writeError(w, http.StatusNotFound, "vehicle_not_found")
	case errors.Is(err, domain.ErrVehicleAlreadyExists):
		writeError(w, http.StatusConflict, "vehicle_already_exists")
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
