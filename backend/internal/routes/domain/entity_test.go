package domain

import "testing"

func TestRouteValidate(t *testing.T) {
	valid := Route{
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   RouteStatusActive,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid route, got %v", err)
	}

	tests := []struct {
		name  string
		route Route
	}{
		{name: "missing company", route: Route{Name: "Bogota - Tunja", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDurationMinutes: 180, Status: RouteStatusActive}},
		{name: "missing name", route: Route{CompanyID: "company-1", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDurationMinutes: 180, Status: RouteStatusActive}},
		{name: "missing origin", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", DestinationCity: "Tunja", EstimatedDurationMinutes: 180, Status: RouteStatusActive}},
		{name: "missing destination", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", OriginCity: "Bogota", EstimatedDurationMinutes: 180, Status: RouteStatusActive}},
		{name: "negative distance", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDistanceKM: -1, EstimatedDurationMinutes: 180, Status: RouteStatusActive}},
		{name: "invalid duration", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDurationMinutes: 0, Status: RouteStatusActive}},
		{name: "negative price", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDurationMinutes: 180, BasePrice: -1, Status: RouteStatusActive}},
		{name: "invalid status", route: Route{CompanyID: "company-1", Name: "Bogota - Tunja", OriginCity: "Bogota", DestinationCity: "Tunja", EstimatedDurationMinutes: 180, Status: "unknown"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.route.Validate(); err != ErrInvalidRoute {
				t.Fatalf("expected ErrInvalidRoute, got %v", err)
			}
		})
	}
}

func TestRouteStatusValid(t *testing.T) {
	for _, status := range []RouteStatus{RouteStatusActive, RouteStatusInactive, RouteStatusArchived} {
		if !status.Valid() {
			t.Fatalf("expected %q to be valid", status)
		}
	}

	if RouteStatus("maintenance").Valid() {
		t.Fatal("expected maintenance to be invalid")
	}
}
