package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type DeactivateDriverUseCase struct {
	drivers domain.Repository
}

func NewDeactivateDriverUseCase(drivers domain.Repository) *DeactivateDriverUseCase {
	return &DeactivateDriverUseCase{drivers: drivers}
}

func (uc *DeactivateDriverUseCase) Execute(ctx context.Context, req dto.DeactivateDriverRequest) error {
	if req.CompanyID == "" || req.ID == "" {
		return domain.ErrInvalidDriver
	}

	current, err := uc.drivers.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return err
	}
	if err := validateDriverStatusTransition(current.Status, domain.DriverStatusInactive); err != nil {
		return err
	}

	return uc.drivers.Deactivate(ctx, req.CompanyID, req.ID)
}
