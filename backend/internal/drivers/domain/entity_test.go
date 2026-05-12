package domain

import "testing"

func TestDriverValidate(t *testing.T) {
	valid := Driver{
		CompanyID:      "11111111-1111-4111-8111-111111111111",
		FirstName:      "Demo",
		LastName:       "Driver",
		DocumentNumber: "DOC123",
		Phone:          "+573001112233",
		Email:          "driver@example.com",
		LicenseNumber:  "LIC123",
		Status:         DriverStatusActive,
	}

	tests := []struct {
		name    string
		mutate  func(*Driver)
		wantErr bool
	}{
		{name: "valid driver"},
		{
			name: "email optional",
			mutate: func(d *Driver) {
				d.Email = ""
			},
		},
		{
			name: "missing company",
			mutate: func(d *Driver) {
				d.CompanyID = ""
			},
			wantErr: true,
		},
		{
			name: "missing document",
			mutate: func(d *Driver) {
				d.DocumentNumber = " "
			},
			wantErr: true,
		},
		{
			name: "blank email when present",
			mutate: func(d *Driver) {
				d.Email = " "
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			mutate: func(d *Driver) {
				d.Status = "deleted"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := valid
			if tt.mutate != nil {
				tt.mutate(&driver)
			}

			err := driver.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestDriverStatusValid(t *testing.T) {
	validStatuses := []DriverStatus{
		DriverStatusActive,
		DriverStatusInactive,
		DriverStatusSuspended,
	}

	for _, status := range validStatuses {
		if !status.Valid() {
			t.Fatalf("expected status %q to be valid", status)
		}
	}

	if DriverStatus("maintenance").Valid() {
		t.Fatal("expected unknown status to be invalid")
	}
}
