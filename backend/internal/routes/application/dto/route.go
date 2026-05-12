package dto

import "time"

type Route struct {
	ID                       string
	CompanyID                string
	Name                     string
	OriginCity               string
	DestinationCity          string
	EstimatedDistanceKM      float64
	EstimatedDurationMinutes int
	BasePrice                float64
	Status                   string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type CreateRouteRequest struct {
	CompanyID                string
	Name                     string
	OriginCity               string
	DestinationCity          string
	EstimatedDistanceKM      float64
	EstimatedDurationMinutes int
	BasePrice                float64
	Status                   string
}

type UpdateRouteRequest struct {
	CompanyID                string
	ID                       string
	Name                     string
	OriginCity               string
	DestinationCity          string
	EstimatedDistanceKM      float64
	EstimatedDurationMinutes int
	BasePrice                float64
	Status                   string
}

type GetRouteRequest struct {
	CompanyID string
	ID        string
}

type ListRoutesRequest struct {
	CompanyID string
}

type ArchiveRouteRequest struct {
	CompanyID string
	ID        string
}
