package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

type ListVehiclesUseCase struct {
	vehicles domain.Repository
}

func NewListVehiclesUseCase(vehicles domain.Repository) *ListVehiclesUseCase {
	return &ListVehiclesUseCase{vehicles: vehicles}
}

func (uc *ListVehiclesUseCase) Execute(ctx context.Context, req dto.ListVehiclesRequest) ([]dto.Vehicle, error) {
	if req.CompanyID == "" {
		return nil, domain.ErrInvalidVehicle
	}

	vehicles, err := uc.vehicles.ListByCompany(ctx, req.CompanyID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.Vehicle, 0, len(vehicles))
	for _, vehicle := range vehicles {
		res = append(res, vehicleToDTO(vehicle))
	}
	return res, nil
}
