package domain

import "testing"

func TestRouteStopValidate(t *testing.T) {
	valid := RouteStop{
		RouteID:   "route-1",
		Name:      "Terminal Norte",
		City:      "Bogota",
		StopOrder: 1,
		Latitude:  4.710989,
		Longitude: -74.072092,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid route stop, got %v", err)
	}

	tests := []struct {
		name string
		stop RouteStop
	}{
		{name: "missing route", stop: RouteStop{Name: "Stop", City: "Bogota", StopOrder: 1}},
		{name: "missing name", stop: RouteStop{RouteID: "route-1", City: "Bogota", StopOrder: 1}},
		{name: "missing city", stop: RouteStop{RouteID: "route-1", Name: "Stop", StopOrder: 1}},
		{name: "invalid order", stop: RouteStop{RouteID: "route-1", Name: "Stop", City: "Bogota", StopOrder: 0}},
		{name: "invalid latitude", stop: RouteStop{RouteID: "route-1", Name: "Stop", City: "Bogota", StopOrder: 1, Latitude: 91}},
		{name: "invalid longitude", stop: RouteStop{RouteID: "route-1", Name: "Stop", City: "Bogota", StopOrder: 1, Longitude: 181}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.stop.Validate(); err != ErrInvalidRouteStop {
				t.Fatalf("expected ErrInvalidRouteStop, got %v", err)
			}
		})
	}
}
