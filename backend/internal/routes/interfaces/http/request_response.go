package http

import "time"

type routeRequest struct {
	Name                     string  `json:"name"`
	OriginCity               string  `json:"origin_city"`
	DestinationCity          string  `json:"destination_city"`
	EstimatedDistanceKM      float64 `json:"estimated_distance_km"`
	EstimatedDurationMinutes int     `json:"estimated_duration_minutes"`
	BasePrice                float64 `json:"base_price"`
	Status                   string  `json:"status"`
}

type routeResponse struct {
	ID                       string    `json:"id"`
	CompanyID                string    `json:"company_id"`
	Name                     string    `json:"name"`
	OriginCity               string    `json:"origin_city"`
	DestinationCity          string    `json:"destination_city"`
	EstimatedDistanceKM      float64   `json:"estimated_distance_km"`
	EstimatedDurationMinutes int       `json:"estimated_duration_minutes"`
	BasePrice                float64   `json:"base_price"`
	Status                   string    `json:"status"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type routeEnvelope struct {
	Route routeResponse `json:"route"`
}

type routesEnvelope struct {
	Routes []routeResponse `json:"routes"`
}

type errorResponse struct {
	Error string `json:"error"`
}
