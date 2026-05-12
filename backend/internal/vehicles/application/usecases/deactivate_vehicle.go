package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type DeactivateVehicleUseCase struct {
	vehicles domain.Repository
}

func NewDeactivateVehicleUseCase(vehicles domain.Repository) *DeactivateVehicleUseCase {
	return &DeactivateVehicleUseCase{vehicles: vehicles}
}

func (uc *DeactivateVehicleUseCase) Execute(ctx context.Context, req dto.DeactivateVehicleRequest) error {
	if req.CompanyID == "" || req.ID == "" {
		return domain.ErrInvalidVehicle
	}

	current, err := uc.vehicles.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return err
	}
	if err := validateStatusTransition(current.Status, domain.VehicleStatusInactive); err != nil {
		return err
	}

	return uc.vehicles.MarkInactive(ctx, req.CompanyID, req.ID)
}
