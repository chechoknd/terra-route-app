package domain

import (
	"strings"
	"time"
)

type RouteStop struct {
	ID        string
	RouteID   string
	Name      string
	City      string
	StopOrder int
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s RouteStop) Validate() error {
	if strings.TrimSpace(s.RouteID) == "" {
		return ErrInvalidRouteStop
	}
	if strings.TrimSpace(s.Name) == "" {
		return ErrInvalidRouteStop
	}
	if strings.TrimSpace(s.City) == "" {
		return ErrInvalidRouteStop
	}
	if s.StopOrder <= 0 {
		return ErrInvalidRouteStop
	}
	if s.Latitude < -90 || s.Latitude > 90 {
		return ErrInvalidRouteStop
	}
	if s.Longitude < -180 || s.Longitude > 180 {
		return ErrInvalidRouteStop
	}
	return nil
}
