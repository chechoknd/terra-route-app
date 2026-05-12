package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

type UpdateDriverUseCase struct {
	drivers domain.Repository
}

func NewUpdateDriverUseCase(drivers domain.Repository) *UpdateDriverUseCase {
	return &UpdateDriverUseCase{drivers: drivers}
}

func (uc *UpdateDriverUseCase) Execute(ctx context.Context, req dto.UpdateDriverRequest) (*dto.Driver, error) {
	if req.ID == "" {
		return nil, domain.ErrInvalidDriver
	}

	current, err := uc.drivers.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	driver := &domain.Driver{
		ID:             current.ID,
		CompanyID:      current.CompanyID,
		UserID:         req.UserID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocumentNumber: req.DocumentNumber,
		Phone:          req.Phone,
		Email:          req.Email,
		LicenseNumber:  req.LicenseNumber,
		Status:         domain.DriverStatus(req.Status),
		CreatedAt:      current.CreatedAt,
		UpdatedAt:      current.UpdatedAt,
	}
	if err := validateDriverStatusTransition(current.Status, driver.Status); err != nil {
		return nil, err
	}
	if err := driver.Validate(); err != nil {
		return nil, err
	}

	if err := uc.drivers.Update(ctx, driver); err != nil {
		return nil, err
	}

	res := driverToDTO(*driver)
	return &res, nil
}

func validateDriverStatusTransition(from domain.DriverStatus, to domain.DriverStatus) error {
	if !to.Valid() {
		return domain.ErrInvalidDriver
	}
	if from == domain.DriverStatusInactive && to != domain.DriverStatusInactive {
		return domain.ErrInvalidDriver
	}
	return nil
}
