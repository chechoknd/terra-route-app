package infrastructure

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	companydomain "github.com/terraroute/terra-route/backend/internal/companies/domain"
	companyinfra "github.com/terraroute/terra-route/backend/internal/companies/infrastructure"
	"github.com/terraroute/terra-route/backend/internal/database"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

func TestPostgresRepositoryVehicleLifecycle(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	companyRepo := companyinfra.NewPostgresRepository(db)
	companyA := createVehicleTestCompany(t, ctx, companyRepo, "a")
	companyB := createVehicleTestCompany(t, ctx, companyRepo, "b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM vehicles WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	vehicle := &domain.Vehicle{
		CompanyID:    companyA.ID,
		Plate:        "ABC-" + safeName(t.Name()),
		InternalCode: "BUS-" + safeName(t.Name()),
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusActive,
	}

	if err := repo.Create(ctx, vehicle); err != nil {
		t.Fatalf("create vehicle: %v", err)
	}

	got, err := repo.GetByID(ctx, companyA.ID, vehicle.ID)
	if err != nil {
		t.Fatalf("get vehicle by id: %v", err)
	}
	if got.ID != vehicle.ID {
		t.Fatalf("expected vehicle id %q, got %q", vehicle.ID, got.ID)
	}

	_, err = repo.GetByID(ctx, companyB.ID, vehicle.ID)
	if !errors.Is(err, domain.ErrVehicleNotFound) {
		t.Fatalf("expected ErrVehicleNotFound for another company scope, got %v", err)
	}

	vehicles, err := repo.ListByCompany(ctx, companyA.ID)
	if err != nil {
		t.Fatalf("list vehicles by company: %v", err)
	}
	if len(vehicles) != 1 {
		t.Fatalf("expected 1 vehicle in company scope, got %d", len(vehicles))
	}

	vehicle.Status = domain.VehicleStatusMaintenance
	vehicle.Capacity = 44
	if err := repo.Update(ctx, vehicle); err != nil {
		t.Fatalf("update vehicle: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, vehicle.ID)
	if err != nil {
		t.Fatalf("get updated vehicle: %v", err)
	}
	if got.Status != domain.VehicleStatusMaintenance || got.Capacity != 44 {
		t.Fatalf("unexpected updated vehicle: %+v", got)
	}

	if err := repo.MarkInactive(ctx, companyA.ID, vehicle.ID); err != nil {
		t.Fatalf("mark inactive: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, vehicle.ID)
	if err != nil {
		t.Fatalf("get inactive vehicle: %v", err)
	}
	if got.Status != domain.VehicleStatusInactive {
		t.Fatalf("expected inactive status, got %q", got.Status)
	}
}

func TestPostgresRepositoryUniquePlatePerCompany(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	companyRepo := companyinfra.NewPostgresRepository(db)
	companyA := createVehicleTestCompany(t, ctx, companyRepo, "unique-a")
	companyB := createVehicleTestCompany(t, ctx, companyRepo, "unique-b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM vehicles WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	plate := "DUP-" + safeName(t.Name())

	first := newTestVehicle(companyA.ID, plate, "A-"+safeName(t.Name()))
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("create first vehicle: %v", err)
	}

	duplicateSameCompany := newTestVehicle(companyA.ID, strings.ToLower(plate), "B-"+safeName(t.Name()))
	err = repo.Create(ctx, duplicateSameCompany)
	if !errors.Is(err, domain.ErrVehicleAlreadyExists) {
		t.Fatalf("expected ErrVehicleAlreadyExists, got %v", err)
	}

	samePlateAnotherCompany := newTestVehicle(companyB.ID, strings.ToLower(plate), "C-"+safeName(t.Name()))
	if err := repo.Create(ctx, samePlateAnotherCompany); err != nil {
		t.Fatalf("same plate in another company should be allowed: %v", err)
	}
}

func createVehicleTestCompany(t *testing.T, ctx context.Context, repo *companyinfra.PostgresRepository, suffix string) *companydomain.Company {
	t.Helper()

	company := &companydomain.Company{
		Name:   "Vehicle Test Company " + suffix,
		Slug:   "vehicle-test-" + safeName(t.Name()) + "-" + suffix,
		Status: companydomain.CompanyStatusActive,
	}
	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company %s: %v", suffix, err)
	}
	return company
}

func newTestVehicle(companyID string, plate string, internalCode string) *domain.Vehicle {
	return &domain.Vehicle{
		CompanyID:    companyID,
		Plate:        plate,
		InternalCode: internalCode,
		VehicleType:  "bus",
		Brand:        "Mercedes-Benz",
		Model:        "OF-1721",
		Capacity:     42,
		Status:       domain.VehicleStatusActive,
	}
}

func safeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
