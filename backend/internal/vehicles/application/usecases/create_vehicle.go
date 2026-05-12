package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type CreateVehicleUseCase struct {
	vehicles domain.Repository
}

func NewCreateVehicleUseCase(vehicles domain.Repository) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{vehicles: vehicles}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, req dto.CreateVehicleRequest) (*dto.Vehicle, error) {
	vehicle := &domain.Vehicle{
		CompanyID:    req.CompanyID,
		Plate:        req.Plate,
		InternalCode: req.InternalCode,
		VehicleType:  req.VehicleType,
		Brand:        req.Brand,
		Model:        req.Model,
		Capacity:     req.Capacity,
		Status:       vehicleStatusOrDefault(req.Status),
	}
	if err := vehicle.Validate(); err != nil {
		return nil, err
	}

	if err := uc.vehicles.Create(ctx, vehicle); err != nil {
		return nil, err
	}

	res := vehicleToDTO(*vehicle)
	return &res, nil
}
