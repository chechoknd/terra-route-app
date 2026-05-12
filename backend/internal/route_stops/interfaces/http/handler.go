package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type AddRouteStopUseCase interface {
	Execute(ctx context.Context, req dto.AddRouteStopRequest) (*dto.RouteStop, error)
}

type UpdateRouteStopUseCase interface {
	Execute(ctx context.Context, req dto.UpdateRouteStopRequest) (*dto.RouteStop, error)
}

type ListRouteStopsUseCase interface {
	Execute(ctx context.Context, req dto.ListRouteStopsRequest) ([]dto.RouteStop, error)
}

type DeleteRouteStopUseCase interface {
	Execute(ctx context.Context, req dto.DeleteRouteStopRequest) error
}

type Handler struct {
	add    AddRouteStopUseCase
	update UpdateRouteStopUseCase
	list   ListRouteStopsUseCase
	delete DeleteRouteStopUseCase
}

func NewHandler(add AddRouteStopUseCase, update UpdateRouteStopUseCase, list ListRouteStopsUseCase, delete DeleteRouteStopUseCase) *Handler {
	return &Handler{add: add, update: update, list: list, delete: delete}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}
	var req routeStopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	stop, err := h.add.Execute(r.Context(), dto.AddRouteStopRequest{
		CompanyID: claims.CompanyID,
		RouteID:   r.PathValue("id"),
		Name:      req.Name,
		City:      req.City,
		StopOrder: req.StopOrder,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		writeRouteStopError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, routeStopEnvelope{Stop: routeStopToResponse(*stop)})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}
	stops, err := h.list.Execute(r.Context(), dto.ListRouteStopsRequest{
		CompanyID: claims.CompanyID,
		RouteID:   r.PathValue("id"),
	})
	if err != nil {
		writeRouteStopError(w, err)
		return
	}
	res := make([]routeStopResponse, 0, len(stops))
	for _, stop := range stops {
		res = append(res, routeStopToResponse(stop))
	}
	writeJSON(w, http.StatusOK, routeStopsEnvelope{Stops: res})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}
	var req routeStopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	stop, err := h.update.Execute(r.Context(), dto.UpdateRouteStopRequest{
		CompanyID: claims.CompanyID,
		RouteID:   r.PathValue("id"),
		ID:        r.PathValue("stopId"),
		Name:      req.Name,
		City:      req.City,
		StopOrder: req.StopOrder,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		writeRouteStopError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, routeStopEnvelope{Stop: routeStopToResponse(*stop)})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := authhttp.AuthClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not_authenticated")
		return
	}
	err := h.delete.Execute(r.Context(), dto.DeleteRouteStopRequest{
		CompanyID: claims.CompanyID,
		RouteID:   r.PathValue("id"),
		ID:        r.PathValue("stopId"),
	})
	if err != nil {
		writeRouteStopError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func routeStopToResponse(stop dto.RouteStop) routeStopResponse {
	return routeStopResponse{
		ID:        stop.ID,
		RouteID:   stop.RouteID,
		Name:      stop.Name,
		City:      stop.City,
		StopOrder: stop.StopOrder,
		Latitude:  stop.Latitude,
		Longitude: stop.Longitude,
		CreatedAt: stop.CreatedAt,
		UpdatedAt: stop.UpdatedAt,
	}
}

func writeRouteStopError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRouteStop):
		writeError(w, http.StatusBadRequest, "invalid_route_stop")
	case errors.Is(err, domain.ErrRouteStopNotFound):
		writeError(w, http.StatusNotFound, "route_stop_not_found")
	case errors.Is(err, domain.ErrRouteStopAlreadyExists):
		writeError(w, http.StatusConflict, "route_stop_already_exists")
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
