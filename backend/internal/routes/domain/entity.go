package domain

import (
	"strings"
	"time"
)

type RouteStatus string

const (
	RouteStatusActive   RouteStatus = "active"
	RouteStatusInactive RouteStatus = "inactive"
	RouteStatusArchived RouteStatus = "archived"
)

type Route struct {
	ID                       string
	CompanyID                string
	Name                     string
	OriginCity               string
	DestinationCity          string
	EstimatedDistanceKM      float64
	EstimatedDurationMinutes int
	BasePrice                float64
	Status                   RouteStatus
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

func (r Route) Validate() error {
	if strings.TrimSpace(r.CompanyID) == "" {
		return ErrInvalidRoute
	}
	if strings.TrimSpace(r.Name) == "" {
		return ErrInvalidRoute
	}
	if strings.TrimSpace(r.OriginCity) == "" {
		return ErrInvalidRoute
	}
	if strings.TrimSpace(r.DestinationCity) == "" {
		return ErrInvalidRoute
	}
	if r.EstimatedDistanceKM < 0 {
		return ErrInvalidRoute
	}
	if r.EstimatedDurationMinutes <= 0 {
		return ErrInvalidRoute
	}
	if r.BasePrice < 0 {
		return ErrInvalidRoute
	}
	if !r.Status.Valid() {
		return ErrInvalidRoute
	}
	return nil
}

func (s RouteStatus) Valid() bool {
	switch s {
	case RouteStatusActive, RouteStatusInactive, RouteStatusArchived:
		return true
	default:
		return false
	}
}
