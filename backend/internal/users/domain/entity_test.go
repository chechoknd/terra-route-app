package domain

import "testing"

func TestUserValidateCompanyScope(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "company user requires company",
			user: User{
				Email:        "admin@example.com",
				FullName:     "Company Admin",
				Role:         UserRoleCompanyAdmin,
				Status:       UserStatusActive,
				PasswordHash: "hash",
			},
			wantErr: true,
		},
		{
			name: "company user with company is valid",
			user: User{
				CompanyID:    "f9a9de1b-1d6b-4e4f-bb7d-f9e8ef0f6d9f",
				Email:        "admin@example.com",
				FullName:     "Company Admin",
				Role:         UserRoleCompanyAdmin,
				Status:       UserStatusActive,
				PasswordHash: "hash",
			},
		},
		{
			name: "super admin must not have company",
			user: User{
				CompanyID:    "f9a9de1b-1d6b-4e4f-bb7d-f9e8ef0f6d9f",
				Email:        "root@example.com",
				FullName:     "Root Admin",
				Role:         UserRoleSuperAdmin,
				Status:       UserStatusActive,
				PasswordHash: "hash",
			},
			wantErr: true,
		},
		{
			name: "super admin without company is valid",
			user: User{
				Email:        "root@example.com",
				FullName:     "Root Admin",
				Role:         UserRoleSuperAdmin,
				Status:       UserStatusActive,
				PasswordHash: "hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
