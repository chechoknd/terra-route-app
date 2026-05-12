package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type UpdateVehicleUseCase struct {
	vehicles domain.Repository
}

func NewUpdateVehicleUseCase(vehicles domain.Repository) *UpdateVehicleUseCase {
	return &UpdateVehicleUseCase{vehicles: vehicles}
}

func (uc *UpdateVehicleUseCase) Execute(ctx context.Context, req dto.UpdateVehicleRequest) (*dto.Vehicle, error) {
	if req.ID == "" {
		return nil, domain.ErrInvalidVehicle
	}

	current, err := uc.vehicles.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	vehicle := &domain.Vehicle{
		ID:           current.ID,
		CompanyID:    current.CompanyID,
		Plate:        req.Plate,
		InternalCode: req.InternalCode,
		VehicleType:  req.VehicleType,
		Brand:        req.Brand,
		Model:        req.Model,
		Capacity:     req.Capacity,
		Status:       domain.VehicleStatus(req.Status),
		CreatedAt:    current.CreatedAt,
		UpdatedAt:    current.UpdatedAt,
	}
	if err := validateStatusTransition(current.Status, vehicle.Status); err != nil {
		return nil, err
	}
	if err := vehicle.Validate(); err != nil {
		return nil, err
	}

	if err := uc.vehicles.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	res := vehicleToDTO(*vehicle)
	return &res, nil
}

func validateStatusTransition(from domain.VehicleStatus, to domain.VehicleStatus) error {
	if !to.Valid() {
		return domain.ErrInvalidVehicle
	}
	if from == domain.VehicleStatusInactive && to != domain.VehicleStatusInactive {
		return domain.ErrInvalidVehicle
	}
	return nil
}
