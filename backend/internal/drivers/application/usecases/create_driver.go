package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type CreateDriverUseCase struct {
	drivers domain.Repository
}

func NewCreateDriverUseCase(drivers domain.Repository) *CreateDriverUseCase {
	return &CreateDriverUseCase{drivers: drivers}
}

func (uc *CreateDriverUseCase) Execute(ctx context.Context, req dto.CreateDriverRequest) (*dto.Driver, error) {
	driver := &domain.Driver{
		CompanyID:      req.CompanyID,
		UserID:         req.UserID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocumentNumber: req.DocumentNumber,
		Phone:          req.Phone,
		Email:          req.Email,
		LicenseNumber:  req.LicenseNumber,
		Status:         driverStatusOrDefault(req.Status),
	}
	if err := driver.Validate(); err != nil {
		return nil, err
	}

	if err := uc.drivers.Create(ctx, driver); err != nil {
		return nil, err
	}

	res := driverToDTO(*driver)
	return &res, nil
}
