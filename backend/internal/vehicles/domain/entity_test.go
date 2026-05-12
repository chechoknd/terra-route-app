package domain

import "testing"

func TestVehicleValidate(t *testing.T) {
	valid := Vehicle{
		CompanyID:    "11111111-1111-4111-8111-111111111111",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       VehicleStatusActive,
	}

	tests := []struct {
		name    string
		mutate  func(*Vehicle)
		wantErr bool
	}{
		{name: "valid vehicle"},
		{
			name: "missing company",
			mutate: func(v *Vehicle) {
				v.CompanyID = ""
			},
			wantErr: true,
		},
		{
			name: "missing plate",
			mutate: func(v *Vehicle) {
				v.Plate = " "
			},
			wantErr: true,
		},
		{
			name: "invalid capacity",
			mutate: func(v *Vehicle) {
				v.Capacity = 0
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			mutate: func(v *Vehicle) {
				v.Status = "retired"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vehicle := valid
			if tt.mutate != nil {
				tt.mutate(&vehicle)
			}

			err := vehicle.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestVehicleStatusValid(t *testing.T) {
	validStatuses := []VehicleStatus{
		VehicleStatusActive,
		VehicleStatusInactive,
		VehicleStatusMaintenance,
		VehicleStatusUnavailable,
	}

	for _, status := range validStatuses {
		if !status.Valid() {
			t.Fatalf("expected status %q to be valid", status)
		}
	}

	if VehicleStatus("deleted").Valid() {
		t.Fatal("expected unknown status to be invalid")
	}
}
