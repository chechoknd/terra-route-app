package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type GetVehicleUseCase struct {
	vehicles domain.Repository
}

func NewGetVehicleUseCase(vehicles domain.Repository) *GetVehicleUseCase {
	return &GetVehicleUseCase{vehicles: vehicles}
}

func (uc *GetVehicleUseCase) Execute(ctx context.Context, req dto.GetVehicleRequest) (*dto.Vehicle, error) {
	if req.CompanyID == "" || req.ID == "" {
		return nil, domain.ErrInvalidVehicle
	}

	vehicle, err := uc.vehicles.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	res := vehicleToDTO(*vehicle)
	return &res, nil
}
