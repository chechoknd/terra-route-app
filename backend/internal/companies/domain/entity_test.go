package domain

import "testing"

func TestCompanyValidate(t *testing.T) {
	tests := []struct {
		name    string
		company Company
		wantErr bool
	}{
		{
			name: "valid company",
			company: Company{
				Name:   "Terra Bus",
				Slug:   "terra-bus",
				Status: CompanyStatusActive,
			},
		},
		{
			name: "missing name",
			company: Company{
				Slug:   "terra-bus",
				Status: CompanyStatusActive,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			company: Company{
				Name:   "Terra Bus",
				Slug:   "terra-bus",
				Status: "deleted",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.company.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
