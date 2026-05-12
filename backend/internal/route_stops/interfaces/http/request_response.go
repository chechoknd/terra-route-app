package http

import "time"

type routeStopRequest struct {
	Name      string  `json:"name"`
	City      string  `json:"city"`
	StopOrder int     `json:"stop_order"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type routeStopResponse struct {
	ID        string    `json:"id"`
	RouteID   string    `json:"route_id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	StopOrder int       `json:"stop_order"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type routeStopEnvelope struct {
	Stop routeStopResponse `json:"stop"`
}

type routeStopsEnvelope struct {
	Stops []routeStopResponse `json:"stops"`
}

type errorResponse struct {
	Error string `json:"error"`
}
