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
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

func TestPostgresRepositoryRouteLifecycle(t *testing.T) {
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
	companyA := createRouteTestCompany(t, ctx, companyRepo, "a")
	companyB := createRouteTestCompany(t, ctx, companyRepo, "b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM routes WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	route := &domain.Route{
		CompanyID:                companyA.ID,
		Name:                     "Route " + safeName(t.Name()),
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   domain.RouteStatusActive,
	}

	if err := repo.Create(ctx, route); err != nil {
		t.Fatalf("create route: %v", err)
	}

	got, err := repo.GetByID(ctx, companyA.ID, route.ID)
	if err != nil {
		t.Fatalf("get route by id: %v", err)
	}
	if got.ID != route.ID {
		t.Fatalf("expected route id %q, got %q", route.ID, got.ID)
	}

	_, err = repo.GetByID(ctx, companyB.ID, route.ID)
	if !errors.Is(err, domain.ErrRouteNotFound) {
		t.Fatalf("expected ErrRouteNotFound for another company scope, got %v", err)
	}

	routes, err := repo.ListByCompany(ctx, companyA.ID)
	if err != nil {
		t.Fatalf("list routes by company: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route in company scope, got %d", len(routes))
	}

	route.Name = "Updated " + route.Name
	route.Status = domain.RouteStatusInactive
	if err := repo.Update(ctx, route); err != nil {
		t.Fatalf("update route: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, route.ID)
	if err != nil {
		t.Fatalf("get updated route: %v", err)
	}
	if got.Status != domain.RouteStatusInactive || got.Name != route.Name {
		t.Fatalf("unexpected updated route: %+v", got)
	}

	if err := repo.Archive(ctx, companyA.ID, route.ID); err != nil {
		t.Fatalf("archive route: %v", err)
	}

	got, err = repo.GetByID(ctx, companyA.ID, route.ID)
	if err != nil {
		t.Fatalf("get archived route: %v", err)
	}
	if got.Status != domain.RouteStatusArchived {
		t.Fatalf("expected archived status, got %q", got.Status)
	}
}

func TestPostgresRepositoryUniqueRouteNamePerCompany(t *testing.T) {
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
	companyA := createRouteTestCompany(t, ctx, companyRepo, "unique-a")
	companyB := createRouteTestCompany(t, ctx, companyRepo, "unique-b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM routes WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	name := "Route Duplicate " + safeName(t.Name())

	first := newTestRoute(companyA.ID, name)
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("create first route: %v", err)
	}

	duplicateSameCompany := newTestRoute(companyA.ID, strings.ToLower(name))
	err = repo.Create(ctx, duplicateSameCompany)
	if !errors.Is(err, domain.ErrRouteAlreadyExists) {
		t.Fatalf("expected ErrRouteAlreadyExists, got %v", err)
	}

	sameNameAnotherCompany := newTestRoute(companyB.ID, strings.ToLower(name))
	if err := repo.Create(ctx, sameNameAnotherCompany); err != nil {
		t.Fatalf("same name in another company should be allowed: %v", err)
	}
}

func createRouteTestCompany(t *testing.T, ctx context.Context, repo *companyinfra.PostgresRepository, suffix string) *companydomain.Company {
	t.Helper()

	company := &companydomain.Company{
		Name:   "Route Test Company " + suffix,
		Slug:   "route-test-" + safeName(t.Name()) + "-" + suffix,
		Status: companydomain.CompanyStatusActive,
	}
	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company %s: %v", suffix, err)
	}
	return company
}

func newTestRoute(companyID string, name string) *domain.Route {
	return &domain.Route{
		CompanyID:                companyID,
		Name:                     name,
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   domain.RouteStatusActive,
	}
}

func safeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
