package dto

import "time"

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

type AddRouteStopRequest struct {
	CompanyID string
	RouteID   string
	Name      string
	City      string
	StopOrder int
	Latitude  float64
	Longitude float64
}

type UpdateRouteStopRequest struct {
	CompanyID string
	RouteID   string
	ID        string
	Name      string
	City      string
	StopOrder int
	Latitude  float64
	Longitude float64
}

type ListRouteStopsRequest struct {
	CompanyID string
	RouteID   string
}

type DeleteRouteStopRequest struct {
	CompanyID string
	RouteID   string
	ID        string
}

type ReorderRouteStopsRequest struct {
	CompanyID  string
	RouteID    string
	OrderedIDs []string
}
