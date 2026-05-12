package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type CreateRouteUseCase interface {
	Execute(ctx context.Context, req dto.CreateRouteRequest) (*dto.Route, error)
}

type UpdateRouteUseCase interface {
	Execute(ctx context.Context, req dto.UpdateRouteRequest) (*dto.Route, error)
}

type GetRouteUseCase interface {
	Execute(ctx context.Context, req dto.GetRouteRequest) (*dto.Route, error)
}

type ListRoutesUseCase interface {
	Execute(ctx context.Context, req dto.ListRoutesRequest) ([]dto.Route, error)
}

type ArchiveRouteUseCase interface {
	Execute(ctx context.Context, req dto.ArchiveRouteRequest) error
}

type Handler struct {
	create  CreateRouteUseCase
	update  UpdateRouteUseCase
	get     GetRouteUseCase
	list    ListRoutesUseCase
	archive ArchiveRouteUseCase
}

func NewHandler(
	create CreateRouteUseCase,
	update UpdateRouteUseCase,
	get GetRouteUseCase,
	list ListRoutesUseCase,
	archive ArchiveRouteUseCase,
) *Handler {
	return &Handler{
		create:  create,
		update:  update,
		get:     get,
		list:    list,
		archive: archive,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	var req routeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	route, err := h.create.Execute(r.Context(), dto.CreateRouteRequest{
		CompanyID:                claims.CompanyID,
		Name:                     req.Name,
		OriginCity:               req.OriginCity,
		DestinationCity:          req.DestinationCity,
		EstimatedDistanceKM:      req.EstimatedDistanceKM,
		EstimatedDurationMinutes: req.EstimatedDurationMinutes,
		BasePrice:                req.BasePrice,
		Status:                   req.Status,
	})
	if err != nil {
		writeRouteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, routeEnvelope{Route: routeToResponse(*route)})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	routes, err := h.list.Execute(r.Context(), dto.ListRoutesRequest{CompanyID: claims.CompanyID})
	if err != nil {
		writeRouteError(w, err)
		return
	}

	res := make([]routeResponse, 0, len(routes))
	for _, route := range routes {
		res = append(res, routeToResponse(route))
	}
	writeJSON(w, http.StatusOK, routesEnvelope{Routes: res})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	route, err := h.get.Execute(r.Context(), dto.GetRouteRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeRouteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, routeEnvelope{Route: routeToResponse(*route)})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	var req routeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	route, err := h.update.Execute(r.Context(), dto.UpdateRouteRequest{
		CompanyID:                claims.CompanyID,
		ID:                       r.PathValue("id"),
		Name:                     req.Name,
		OriginCity:               req.OriginCity,
		DestinationCity:          req.DestinationCity,
		EstimatedDistanceKM:      req.EstimatedDistanceKM,
		EstimatedDurationMinutes: req.EstimatedDurationMinutes,
		BasePrice:                req.BasePrice,
		Status:                   req.Status,
	})
	if err != nil {
		writeRouteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, routeEnvelope{Route: routeToResponse(*route)})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}

	err := h.archive.Execute(r.Context(), dto.ArchiveRouteRequest{
		CompanyID: claims.CompanyID,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeRouteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func routeToResponse(route dto.Route) routeResponse {
	return routeResponse{
		ID:                       route.ID,
		CompanyID:                route.CompanyID,
		Name:                     route.Name,
		OriginCity:               route.OriginCity,
		DestinationCity:          route.DestinationCity,
		EstimatedDistanceKM:      route.EstimatedDistanceKM,
		EstimatedDurationMinutes: route.EstimatedDurationMinutes,
		BasePrice:                route.BasePrice,
		Status:                   route.Status,
		CreatedAt:                route.CreatedAt,
		UpdatedAt:                route.UpdatedAt,
	}
}

func writeRouteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRoute):
		writeError(w, http.StatusBadRequest, "invalid_route")
	case errors.Is(err, domain.ErrRouteNotFound):
		writeError(w, http.StatusNotFound, "route_not_found")
	case errors.Is(err, domain.ErrRouteAlreadyExists):
		writeError(w, http.StatusConflict, "route_already_exists")
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
