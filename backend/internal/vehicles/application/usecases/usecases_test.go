package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

func TestCreateVehicleUseCase(t *testing.T) {
	repo := newFakeVehicleRepository()
	uc := NewCreateVehicleUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.CreateVehicleRequest{
		CompanyID:    "company-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
	})
	if err != nil {
		t.Fatalf("create vehicle: %v", err)
	}

	if res.ID == "" {
		t.Fatal("expected created vehicle id")
	}
	if res.Status != string(domain.VehicleStatusActive) {
		t.Fatalf("expected default active status, got %q", res.Status)
	}
	if repo.created == nil || repo.created.CompanyID != "company-1" {
		t.Fatalf("vehicle was not passed to repository")
	}
}

func TestCreateVehicleUseCaseValidatesRequiredFields(t *testing.T) {
	uc := NewCreateVehicleUseCase(newFakeVehicleRepository())

	_, err := uc.Execute(context.Background(), dto.CreateVehicleRequest{
		CompanyID: "company-1",
	})
	if !errors.Is(err, domain.ErrInvalidVehicle) {
		t.Fatalf("expected ErrInvalidVehicle, got %v", err)
	}
}

func TestGetVehicleUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeVehicleRepository()
	repo.vehicle = &domain.Vehicle{
		ID:           "vehicle-1",
		CompanyID:    "company-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusActive,
	}
	uc := NewGetVehicleUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.GetVehicleRequest{
		CompanyID: "company-1",
		ID:        "vehicle-1",
	})
	if err != nil {
		t.Fatalf("get vehicle: %v", err)
	}

	if res.ID != "vehicle-1" {
		t.Fatalf("expected vehicle-1, got %q", res.ID)
	}
	if repo.getCompanyID != "company-1" || repo.getID != "vehicle-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestListVehiclesUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeVehicleRepository()
	repo.vehicles = []domain.Vehicle{
		{
			ID:           "vehicle-1",
			CompanyID:    "company-1",
			Plate:        "ABC123",
			InternalCode: "BUS-001",
			VehicleType:  "bus",
			Brand:        "Mercedes-Benz",
			Model:        "OF-1721",
			Capacity:     42,
			Status:       domain.VehicleStatusActive,
		},
	}
	uc := NewListVehiclesUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.ListVehiclesRequest{CompanyID: "company-1"})
	if err != nil {
		t.Fatalf("list vehicles: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 vehicle, got %d", len(res))
	}
	if repo.listCompanyID != "company-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestUpdateVehicleUseCase(t *testing.T) {
	repo := newFakeVehicleRepository()
	repo.vehicle = &domain.Vehicle{
		ID:           "vehicle-1",
		CompanyID:    "company-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusActive,
	}
	uc := NewUpdateVehicleUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.UpdateVehicleRequest{
		CompanyID:    "company-1",
		ID:           "vehicle-1",
		Plate:        "XYZ987",
		InternalCode: "BUS-009",
		VehicleType:  "van",
		Brand:        "Renault",
		Model:        "Master",
		Capacity:     18,
		Status:       string(domain.VehicleStatusMaintenance),
	})
	if err != nil {
		t.Fatalf("update vehicle: %v", err)
	}

	if res.Plate != "XYZ987" || res.Status != string(domain.VehicleStatusMaintenance) {
		t.Fatalf("unexpected update response: %+v", res)
	}
	if repo.updated == nil || repo.updated.CompanyID != "company-1" || repo.updated.ID != "vehicle-1" {
		t.Fatalf("repository update was not company scoped")
	}
}

func TestUpdateVehicleUseCaseRejectsInactiveReactivation(t *testing.T) {
	repo := newFakeVehicleRepository()
	repo.vehicle = &domain.Vehicle{
		ID:           "vehicle-1",
		CompanyID:    "company-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusInactive,
	}
	uc := NewUpdateVehicleUseCase(repo)

	_, err := uc.Execute(context.Background(), dto.UpdateVehicleRequest{
		CompanyID:    "company-1",
		ID:           "vehicle-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       string(domain.VehicleStatusActive),
	})
	if !errors.Is(err, domain.ErrInvalidVehicle) {
		t.Fatalf("expected ErrInvalidVehicle, got %v", err)
	}
}

func TestDeactivateVehicleUseCase(t *testing.T) {
	repo := newFakeVehicleRepository()
	repo.vehicle = &domain.Vehicle{
		ID:           "vehicle-1",
		CompanyID:    "company-1",
		Plate:        "ABC123",
		InternalCode: "BUS-001",
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusActive,
	}
	uc := NewDeactivateVehicleUseCase(repo)

	if err := uc.Execute(context.Background(), dto.DeactivateVehicleRequest{
		CompanyID: "company-1",
		ID:        "vehicle-1",
	}); err != nil {
		t.Fatalf("deactivate vehicle: %v", err)
	}

	if repo.inactiveCompanyID != "company-1" || repo.inactiveID != "vehicle-1" {
		t.Fatalf("repository deactivate was not company scoped")
	}
}

type fakeVehicleRepository struct {
	created *domain.Vehicle
	updated *domain.Vehicle

	vehicle  *domain.Vehicle
	vehicles []domain.Vehicle
	err      error

	getCompanyID      string
	getID             string
	listCompanyID     string
	inactiveCompanyID string
	inactiveID        string
}

func newFakeVehicleRepository() *fakeVehicleRepository {
	return &fakeVehicleRepository{}
}

func (r *fakeVehicleRepository) Create(_ context.Context, vehicle *domain.Vehicle) error {
	if r.err != nil {
		return r.err
	}
	vehicle.ID = "vehicle-1"
	r.created = vehicle
	return nil
}

func (r *fakeVehicleRepository) GetByID(_ context.Context, companyID string, id string) (*domain.Vehicle, error) {
	r.getCompanyID = companyID
	r.getID = id
	if r.err != nil {
		return nil, r.err
	}
	if r.vehicle == nil {
		return nil, domain.ErrVehicleNotFound
	}
	return r.vehicle, nil
}

func (r *fakeVehicleRepository) ListByCompany(_ context.Context, companyID string) ([]domain.Vehicle, error) {
	r.listCompanyID = companyID
	if r.err != nil {
		return nil, r.err
	}
	return r.vehicles, nil
}

func (r *fakeVehicleRepository) Update(_ context.Context, vehicle *domain.Vehicle) error {
	if r.err != nil {
		return r.err
	}
	r.updated = vehicle
	return nil
}

func (r *fakeVehicleRepository) MarkInactive(_ context.Context, companyID string, id string) error {
	r.inactiveCompanyID = companyID
	r.inactiveID = id
	return r.err
}
