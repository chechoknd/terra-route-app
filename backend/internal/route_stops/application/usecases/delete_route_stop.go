package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type DeleteRouteStopUseCase struct {
	stops domain.Repository
}

func NewDeleteRouteStopUseCase(stops domain.Repository) *DeleteRouteStopUseCase {
	return &DeleteRouteStopUseCase{stops: stops}
}

func (uc *DeleteRouteStopUseCase) Execute(ctx context.Context, req dto.DeleteRouteStopRequest) error {
	if req.CompanyID == "" || req.RouteID == "" || req.ID == "" {
		return domain.ErrInvalidRouteStop
	}
	return uc.stops.Delete(ctx, req.CompanyID, req.RouteID, req.ID)
}
