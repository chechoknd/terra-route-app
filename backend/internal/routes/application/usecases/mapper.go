package usecases

import (
	"strings"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

func routeToDTO(route domain.Route) dto.Route {
	return dto.Route{
		ID:                       route.ID,
		CompanyID:                route.CompanyID,
		Name:                     route.Name,
		OriginCity:               route.OriginCity,
		DestinationCity:          route.DestinationCity,
		EstimatedDistanceKM:      route.EstimatedDistanceKM,
		EstimatedDurationMinutes: route.EstimatedDurationMinutes,
		BasePrice:                route.BasePrice,
		Status:                   string(route.Status),
		CreatedAt:                route.CreatedAt,
		UpdatedAt:                route.UpdatedAt,
	}
}

func routeStatusOrDefault(status string) domain.RouteStatus {
	if strings.TrimSpace(status) == "" {
		return domain.RouteStatusActive
	}
	return domain.RouteStatus(status)
}
