package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type GetDriverUseCase struct {
	drivers domain.Repository
}

func NewGetDriverUseCase(drivers domain.Repository) *GetDriverUseCase {
	return &GetDriverUseCase{drivers: drivers}
}

func (uc *GetDriverUseCase) Execute(ctx context.Context, req dto.GetDriverRequest) (*dto.Driver, error) {
	if req.CompanyID == "" || req.ID == "" {
		return nil, domain.ErrInvalidDriver
	}

	driver, err := uc.drivers.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	res := driverToDTO(*driver)
	return &res, nil
}
