package usecases

import (
	"strings"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

func vehicleToDTO(vehicle domain.Vehicle) dto.Vehicle {
	return dto.Vehicle{
		ID:           vehicle.ID,
		CompanyID:    vehicle.CompanyID,
		Plate:        vehicle.Plate,
		InternalCode: vehicle.InternalCode,
		VehicleType:  vehicle.VehicleType,
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Capacity:     vehicle.Capacity,
		Status:       string(vehicle.Status),
		CreatedAt:    vehicle.CreatedAt,
		UpdatedAt:    vehicle.UpdatedAt,
	}
}

func vehicleStatusOrDefault(status string) domain.VehicleStatus {
	if strings.TrimSpace(status) == "" {
		return domain.VehicleStatusActive
	}
	return domain.VehicleStatus(status)
}
