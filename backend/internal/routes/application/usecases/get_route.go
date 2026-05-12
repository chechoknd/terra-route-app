package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type GetRouteUseCase struct {
	routes domain.Repository
}

func NewGetRouteUseCase(routes domain.Repository) *GetRouteUseCase {
	return &GetRouteUseCase{routes: routes}
}

func (uc *GetRouteUseCase) Execute(ctx context.Context, req dto.GetRouteRequest) (*dto.Route, error) {
	if req.CompanyID == "" || req.ID == "" {
		return nil, domain.ErrInvalidRoute
	}

	route, err := uc.routes.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	res := routeToDTO(*route)
	return &res, nil
}
