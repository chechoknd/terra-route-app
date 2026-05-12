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
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

func TestPostgresRepositoryDriverLifecycle(t *testing.T) {
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
	companyA := createDriverTestCompany(t, ctx, companyRepo, "a")
	companyB := createDriverTestCompany(t, ctx, companyRepo, "b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM drivers WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	driver := &domain.Driver{
		CompanyID:      companyA.ID,
		FirstName:      "Demo",
		LastName:       "Driver",
		DocumentNumber: "DOC-" + safeName(t.Name()),
		Phone:          "+573001112233",
		Email:          "driver-" + safeName(t.Name()) + "@example.com",
		LicenseNumber:  "LIC-" + safeName(t.Name()),
		Status:         domain.DriverStatusActive,
	}

	if err := repo.Create(ctx, driver); err != nil {
		t.Fatalf("create driver: %v", err)
	}

	got, err := repo.GetByID(ctx, companyA.ID, driver.ID)
	if err != nil {
		t.Fatalf("get driver by id: %v", err)
	}
	if got.ID != driver.ID {
		t.Fatalf("expected driver id %q, got %q", driver.ID, got.ID)
	}

	_, err = repo.GetByID(ctx, companyB.ID, driver.ID)
	if !errors.Is(err, domain.ErrDriverNotFound) {
		t.Fatalf("expected ErrDriverNotFound for another company scope, got %v", err)
	}

	drivers, err := repo.ListByCompany(ctx, companyA.ID)
	if err != nil {
		t.Fatalf("list drivers by company: %v", err)
	}
	if len(drivers) != 1 {
		t.Fatalf("expected 1 driver in company scope, got %d", len(drivers))
	}

	driver.Phone = "+573009998877"
	driver.Status = domain.DriverStatusSuspended
	if err := repo.Update(ctx, driver); err != nil {
		t.Fatalf("update driver: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, driver.ID)
	if err != nil {
		t.Fatalf("get updated driver: %v", err)
	}
	if got.Status != domain.DriverStatusSuspended || got.Phone != "+573009998877" {
		t.Fatalf("unexpected updated driver: %+v", got)
	}

	if err := repo.Deactivate(ctx, companyA.ID, driver.ID); err != nil {
		t.Fatalf("deactivate driver: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, driver.ID)
	if err != nil {
		t.Fatalf("get inactive driver: %v", err)
	}
	if got.Status != domain.DriverStatusInactive {
		t.Fatalf("expected inactive status, got %q", got.Status)
	}
}

func TestPostgresRepositoryUniqueDocumentNumberPerCompany(t *testing.T) {
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
	companyA := createDriverTestCompany(t, ctx, companyRepo, "unique-a")
	companyB := createDriverTestCompany(t, ctx, companyRepo, "unique-b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM drivers WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	document := "DOC-DUP-" + safeName(t.Name())

	first := newTestDriver(companyA.ID, document, "a-"+safeName(t.Name()))
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("create first driver: %v", err)
	}

	duplicateSameCompany := newTestDriver(companyA.ID, strings.ToLower(document), "b-"+safeName(t.Name()))
	err = repo.Create(ctx, duplicateSameCompany)
	if !errors.Is(err, domain.ErrDriverAlreadyExists) {
		t.Fatalf("expected ErrDriverAlreadyExists, got %v", err)
	}

	sameDocumentAnotherCompany := newTestDriver(companyB.ID, strings.ToLower(document), "c-"+safeName(t.Name()))
	if err := repo.Create(ctx, sameDocumentAnotherCompany); err != nil {
		t.Fatalf("same document in another company should be allowed: %v", err)
	}
}

func createDriverTestCompany(t *testing.T, ctx context.Context, repo *companyinfra.PostgresRepository, suffix string) *companydomain.Company {
	t.Helper()

	company := &companydomain.Company{
		Name:   "Driver Test Company " + suffix,
		Slug:   "driver-test-" + safeName(t.Name()) + "-" + suffix,
		Status: companydomain.CompanyStatusActive,
	}
	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company %s: %v", suffix, err)
	}
	return company
}

func newTestDriver(companyID string, documentNumber string, suffix string) *domain.Driver {
	return &domain.Driver{
		CompanyID:      companyID,
		FirstName:      "Demo",
		LastName:       "Driver",
		DocumentNumber: documentNumber,
		Phone:          "+573001112233",
		Email:          "driver-" + suffix + "@example.com",
		LicenseNumber:  "LIC-" + suffix,
		Status:         domain.DriverStatusActive,
	}
}

func safeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
