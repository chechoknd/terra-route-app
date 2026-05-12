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
	stopdomain "github.com/terraroute/terra-route/backend/internal/route_stops/domain"
	routedomain "github.com/terraroute/terra-route/backend/internal/routes/domain"
	routeinfra "github.com/terraroute/terra-route/backend/internal/routes/infrastructure"
)

func TestPostgresRepositoryRouteStopLifecycleAndCompanyScope(t *testing.T) {
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
	companyA := createRouteStopTestCompany(t, ctx, companyRepo, "a")
	companyB := createRouteStopTestCompany(t, ctx, companyRepo, "b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM route_stops WHERE route_id IN (SELECT id FROM routes WHERE company_id IN ($1, $2))", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM routes WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	routeRepo := routeinfra.NewPostgresRepository(db)
	route := createRouteStopTestRoute(t, ctx, routeRepo, companyA.ID, "a")

	repo := NewPostgresRepository(db)
	stop := &stopdomain.RouteStop{
		RouteID:   route.ID,
		Name:      "Terminal Norte",
		City:      "Bogota",
		StopOrder: 1,
		Latitude:  4.710989,
		Longitude: -74.072092,
	}
	if err := repo.Create(ctx, companyA.ID, stop); err != nil {
		t.Fatalf("create route stop: %v", err)
	}

	stops, err := repo.ListByRoute(ctx, companyA.ID, route.ID)
	if err != nil {
		t.Fatalf("list route stops: %v", err)
	}
	if len(stops) != 1 || stops[0].ID != stop.ID {
		t.Fatalf("unexpected route stops: %+v", stops)
	}

	_, err = repo.ListByRoute(ctx, companyB.ID, route.ID)
	if !errors.Is(err, stopdomain.ErrRouteStopNotFound) {
		t.Fatalf("expected ErrRouteStopNotFound for another company scope, got %v", err)
	}

	stop.Name = "Terminal Salitre"
	stop.StopOrder = 2
	if err := repo.Update(ctx, companyA.ID, stop); err != nil {
		t.Fatalf("update route stop: %v", err)
	}

	got, err := repo.GetByID(ctx, companyA.ID, route.ID, stop.ID)
	if err != nil {
		t.Fatalf("get route stop: %v", err)
	}
	if got.Name != "Terminal Salitre" || got.StopOrder != 2 {
		t.Fatalf("unexpected updated stop: %+v", got)
	}

	if err := repo.Delete(ctx, companyA.ID, route.ID, stop.ID); err != nil {
		t.Fatalf("delete route stop: %v", err)
	}
}

func createRouteStopTestCompany(t *testing.T, ctx context.Context, repo *companyinfra.PostgresRepository, suffix string) *companydomain.Company {
	t.Helper()

	company := &companydomain.Company{
		Name:   "Route Stop Test Company " + suffix,
		Slug:   "route-stop-test-" + safeName(t.Name()) + "-" + suffix,
		Status: companydomain.CompanyStatusActive,
	}
	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company %s: %v", suffix, err)
	}
	return company
}

func createRouteStopTestRoute(t *testing.T, ctx context.Context, repo *routeinfra.PostgresRepository, companyID string, suffix string) *routedomain.Route {
	t.Helper()

	route := &routedomain.Route{
		CompanyID:                companyID,
		Name:                     "Route Stop Test Route " + safeName(t.Name()) + "-" + suffix,
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   routedomain.RouteStatusActive,
	}
	if err := repo.Create(ctx, route); err != nil {
		t.Fatalf("create route %s: %v", suffix, err)
	}
	return route
}

func safeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
