package usecases

import (
	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

func routeStopToDTO(stop domain.RouteStop) dto.RouteStop {
	return dto.RouteStop{
		ID:        stop.ID,
		RouteID:   stop.RouteID,
		Name:      stop.Name,
		City:      stop.City,
		StopOrder: stop.StopOrder,
		Latitude:  stop.Latitude,
		Longitude: stop.Longitude,
		CreatedAt: stop.CreatedAt,
		UpdatedAt: stop.UpdatedAt,
	}
}
