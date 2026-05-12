package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type ArchiveRouteUseCase struct {
	routes domain.Repository
}

func NewArchiveRouteUseCase(routes domain.Repository) *ArchiveRouteUseCase {
	return &ArchiveRouteUseCase{routes: routes}
}

func (uc *ArchiveRouteUseCase) Execute(ctx context.Context, req dto.ArchiveRouteRequest) error {
	if req.CompanyID == "" || req.ID == "" {
		return domain.ErrInvalidRoute
	}

	current, err := uc.routes.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return err
	}
	if err := validateRouteStatusTransition(current.Status, domain.RouteStatusArchived); err != nil {
		return err
	}

	return uc.routes.Archive(ctx, req.CompanyID, req.ID)
}
