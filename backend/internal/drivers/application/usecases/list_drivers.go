package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type ListDriversUseCase struct {
	drivers domain.Repository
}

func NewListDriversUseCase(drivers domain.Repository) *ListDriversUseCase {
	return &ListDriversUseCase{drivers: drivers}
}

func (uc *ListDriversUseCase) Execute(ctx context.Context, req dto.ListDriversRequest) ([]dto.Driver, error) {
	if req.CompanyID == "" {
		return nil, domain.ErrInvalidDriver
	}

	drivers, err := uc.drivers.ListByCompany(ctx, req.CompanyID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.Driver, 0, len(drivers))
	for _, driver := range drivers {
		res = append(res, driverToDTO(driver))
	}
	return res, nil
}
