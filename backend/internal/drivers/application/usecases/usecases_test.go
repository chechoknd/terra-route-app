package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

func TestCreateDriverUseCase(t *testing.T) {
	repo := newFakeDriverRepository()
	uc := NewCreateDriverUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.CreateDriverRequest{
		CompanyID:      "company-1",
		UserID:         "user-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		Email:          "ana@example.test",
		LicenseNumber:  "LIC-001",
	})
	if err != nil {
		t.Fatalf("create driver: %v", err)
	}

	if res.ID == "" {
		t.Fatal("expected created driver id")
	}
	if res.Status != string(domain.DriverStatusActive) {
		t.Fatalf("expected default active status, got %q", res.Status)
	}
	if repo.created == nil || repo.created.CompanyID != "company-1" {
		t.Fatalf("driver was not passed to repository")
	}
}

func TestCreateDriverUseCaseValidatesRequiredFields(t *testing.T) {
	uc := NewCreateDriverUseCase(newFakeDriverRepository())

	_, err := uc.Execute(context.Background(), dto.CreateDriverRequest{
		CompanyID: "company-1",
	})
	if !errors.Is(err, domain.ErrInvalidDriver) {
		t.Fatalf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestGetDriverUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeDriverRepository()
	repo.driver = &domain.Driver{
		ID:             "driver-1",
		CompanyID:      "company-1",
		UserID:         "user-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		Email:          "ana@example.test",
		LicenseNumber:  "LIC-001",
		Status:         domain.DriverStatusActive,
	}
	uc := NewGetDriverUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.GetDriverRequest{
		CompanyID: "company-1",
		ID:        "driver-1",
	})
	if err != nil {
		t.Fatalf("get driver: %v", err)
	}

	if res.ID != "driver-1" {
		t.Fatalf("expected driver-1, got %q", res.ID)
	}
	if repo.getCompanyID != "company-1" || repo.getID != "driver-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestListDriversUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeDriverRepository()
	repo.drivers = []domain.Driver{
		{
			ID:             "driver-1",
			CompanyID:      "company-1",
			FirstName:      "Ana",
			LastName:       "Torres",
			DocumentNumber: "DOC-001",
			Phone:          "+573001112233",
			LicenseNumber:  "LIC-001",
			Status:         domain.DriverStatusActive,
		},
	}
	uc := NewListDriversUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.ListDriversRequest{CompanyID: "company-1"})
	if err != nil {
		t.Fatalf("list drivers: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(res))
	}
	if repo.listCompanyID != "company-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestUpdateDriverUseCase(t *testing.T) {
	repo := newFakeDriverRepository()
	repo.driver = &domain.Driver{
		ID:             "driver-1",
		CompanyID:      "company-1",
		UserID:         "user-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		Email:          "ana@example.test",
		LicenseNumber:  "LIC-001",
		Status:         domain.DriverStatusActive,
	}
	uc := NewUpdateDriverUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.UpdateDriverRequest{
		CompanyID:      "company-1",
		ID:             "driver-1",
		UserID:         "user-2",
		FirstName:      "Ana Maria",
		LastName:       "Torres",
		DocumentNumber: "DOC-009",
		Phone:          "+573004445566",
		Email:          "ana.maria@example.test",
		LicenseNumber:  "LIC-009",
		Status:         string(domain.DriverStatusSuspended),
	})
	if err != nil {
		t.Fatalf("update driver: %v", err)
	}

	if res.FirstName != "Ana Maria" || res.Status != string(domain.DriverStatusSuspended) {
		t.Fatalf("unexpected update response: %+v", res)
	}
	if repo.updated == nil || repo.updated.CompanyID != "company-1" || repo.updated.ID != "driver-1" {
		t.Fatalf("repository update was not company scoped")
	}
}

func TestUpdateDriverUseCaseRejectsInactiveReactivation(t *testing.T) {
	repo := newFakeDriverRepository()
	repo.driver = &domain.Driver{
		ID:             "driver-1",
		CompanyID:      "company-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		LicenseNumber:  "LIC-001",
		Status:         domain.DriverStatusInactive,
	}
	uc := NewUpdateDriverUseCase(repo)

	_, err := uc.Execute(context.Background(), dto.UpdateDriverRequest{
		CompanyID:      "company-1",
		ID:             "driver-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		LicenseNumber:  "LIC-001",
		Status:         string(domain.DriverStatusActive),
	})
	if !errors.Is(err, domain.ErrInvalidDriver) {
		t.Fatalf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestDeactivateDriverUseCase(t *testing.T) {
	repo := newFakeDriverRepository()
	repo.driver = &domain.Driver{
		ID:             "driver-1",
		CompanyID:      "company-1",
		FirstName:      "Ana",
		LastName:       "Torres",
		DocumentNumber: "DOC-001",
		Phone:          "+573001112233",
		LicenseNumber:  "LIC-001",
		Status:         domain.DriverStatusActive,
	}
	uc := NewDeactivateDriverUseCase(repo)

	if err := uc.Execute(context.Background(), dto.DeactivateDriverRequest{
		CompanyID: "company-1",
		ID:        "driver-1",
	}); err != nil {
		t.Fatalf("deactivate driver: %v", err)
	}

	if repo.deactivateCompanyID != "company-1" || repo.deactivateID != "driver-1" {
		t.Fatalf("repository deactivate was not company scoped")
	}
}

type fakeDriverRepository struct {
	created *domain.Driver
	updated *domain.Driver

	driver  *domain.Driver
	drivers []domain.Driver
	err     error

	getCompanyID        string
	getID               string
	listCompanyID       string
	deactivateCompanyID string
	deactivateID        string
}

func newFakeDriverRepository() *fakeDriverRepository {
	return &fakeDriverRepository{}
}

func (r *fakeDriverRepository) Create(_ context.Context, driver *domain.Driver) error {
	if r.err != nil {
		return r.err
	}
	driver.ID = "driver-1"
	r.created = driver
	return nil
}

func (r *fakeDriverRepository) GetByID(_ context.Context, companyID string, id string) (*domain.Driver, error) {
	r.getCompanyID = companyID
	r.getID = id
	if r.err != nil {
		return nil, r.err
	}
	if r.driver == nil {
		return nil, domain.ErrDriverNotFound
	}
	return r.driver, nil
}

func (r *fakeDriverRepository) ListByCompany(_ context.Context, companyID string) ([]domain.Driver, error) {
	r.listCompanyID = companyID
	if r.err != nil {
		return nil, r.err
	}
	return r.drivers, nil
}

func (r *fakeDriverRepository) Update(_ context.Context, driver *domain.Driver) error {
	if r.err != nil {
		return r.err
	}
	r.updated = driver
	return nil
}

func (r *fakeDriverRepository) Deactivate(_ context.Context, companyID string, id string) error {
	r.deactivateCompanyID = companyID
	r.deactivateID = id
	return r.err
}
