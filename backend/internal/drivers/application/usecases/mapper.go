package usecases

import (
	"strings"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

func driverToDTO(driver domain.Driver) dto.Driver {
	return dto.Driver{
		ID:             driver.ID,
		CompanyID:      driver.CompanyID,
		UserID:         driver.UserID,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		DocumentNumber: driver.DocumentNumber,
		Phone:          driver.Phone,
		Email:          driver.Email,
		LicenseNumber:  driver.LicenseNumber,
		Status:         string(driver.Status),
		CreatedAt:      driver.CreatedAt,
		UpdatedAt:      driver.UpdatedAt,
	}
}

func driverStatusOrDefault(status string) domain.DriverStatus {
	if strings.TrimSpace(status) == "" {
		return domain.DriverStatusActive
	}
	return domain.DriverStatus(status)
}
