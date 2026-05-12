package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type ListRoutesUseCase struct {
	routes domain.Repository
}

func NewListRoutesUseCase(routes domain.Repository) *ListRoutesUseCase {
	return &ListRoutesUseCase{routes: routes}
}

func (uc *ListRoutesUseCase) Execute(ctx context.Context, req dto.ListRoutesRequest) ([]dto.Route, error) {
	if req.CompanyID == "" {
		return nil, domain.ErrInvalidRoute
	}

	routes, err := uc.routes.ListByCompany(ctx, req.CompanyID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.Route, 0, len(routes))
	for _, route := range routes {
		res = append(res, routeToDTO(route))
	}
	return res, nil
}
