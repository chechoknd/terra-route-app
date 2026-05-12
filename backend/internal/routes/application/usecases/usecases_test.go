package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

func TestCreateRouteUseCase(t *testing.T) {
	repo := newFakeRouteRepository()
	uc := NewCreateRouteUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.CreateRouteRequest{
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
	})
	if err != nil {
		t.Fatalf("create route: %v", err)
	}

	if res.ID == "" {
		t.Fatal("expected created route id")
	}
	if res.Status != string(domain.RouteStatusActive) {
		t.Fatalf("expected default active status, got %q", res.Status)
	}
	if repo.created == nil || repo.created.CompanyID != "company-1" {
		t.Fatalf("route was not passed to repository")
	}
}

func TestCreateRouteUseCaseValidatesRequiredFields(t *testing.T) {
	uc := NewCreateRouteUseCase(newFakeRouteRepository())

	_, err := uc.Execute(context.Background(), dto.CreateRouteRequest{
		CompanyID: "company-1",
	})
	if !errors.Is(err, domain.ErrInvalidRoute) {
		t.Fatalf("expected ErrInvalidRoute, got %v", err)
	}
}

func TestGetRouteUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeRouteRepository()
	repo.route = &domain.Route{
		ID:                       "route-1",
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   domain.RouteStatusActive,
	}
	uc := NewGetRouteUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.GetRouteRequest{
		CompanyID: "company-1",
		ID:        "route-1",
	})
	if err != nil {
		t.Fatalf("get route: %v", err)
	}

	if res.ID != "route-1" {
		t.Fatalf("expected route-1, got %q", res.ID)
	}
	if repo.getCompanyID != "company-1" || repo.getID != "route-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestListRoutesUseCaseScopesByCompany(t *testing.T) {
	repo := newFakeRouteRepository()
	repo.routes = []domain.Route{
		{
			ID:                       "route-1",
			CompanyID:                "company-1",
			Name:                     "Bogota - Tunja",
			OriginCity:               "Bogota",
			DestinationCity:          "Tunja",
			EstimatedDurationMinutes: 180,
			Status:                   domain.RouteStatusActive,
		},
	}
	uc := NewListRoutesUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.ListRoutesRequest{CompanyID: "company-1"})
	if err != nil {
		t.Fatalf("list routes: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 route, got %d", len(res))
	}
	if repo.listCompanyID != "company-1" {
		t.Fatalf("repository was not called with company scope")
	}
}

func TestUpdateRouteUseCase(t *testing.T) {
	repo := newFakeRouteRepository()
	repo.route = &domain.Route{
		ID:                       "route-1",
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDistanceKM:      140.5,
		EstimatedDurationMinutes: 180,
		BasePrice:                45000,
		Status:                   domain.RouteStatusActive,
	}
	uc := NewUpdateRouteUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.UpdateRouteRequest{
		CompanyID:                "company-1",
		ID:                       "route-1",
		Name:                     "Bogota - Duitama",
		OriginCity:               "Bogota",
		DestinationCity:          "Duitama",
		EstimatedDistanceKM:      190.2,
		EstimatedDurationMinutes: 240,
		BasePrice:                60000,
		Status:                   string(domain.RouteStatusInactive),
	})
	if err != nil {
		t.Fatalf("update route: %v", err)
	}

	if res.Name != "Bogota - Duitama" || res.Status != string(domain.RouteStatusInactive) {
		t.Fatalf("unexpected update response: %+v", res)
	}
	if repo.updated == nil || repo.updated.CompanyID != "company-1" || repo.updated.ID != "route-1" {
		t.Fatalf("repository update was not company scoped")
	}
}

func TestUpdateRouteUseCaseRejectsArchivedReactivation(t *testing.T) {
	repo := newFakeRouteRepository()
	repo.route = &domain.Route{
		ID:                       "route-1",
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDurationMinutes: 180,
		Status:                   domain.RouteStatusArchived,
	}
	uc := NewUpdateRouteUseCase(repo)

	_, err := uc.Execute(context.Background(), dto.UpdateRouteRequest{
		CompanyID:                "company-1",
		ID:                       "route-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDurationMinutes: 180,
		Status:                   string(domain.RouteStatusActive),
	})
	if !errors.Is(err, domain.ErrInvalidRoute) {
		t.Fatalf("expected ErrInvalidRoute, got %v", err)
	}
}

func TestArchiveRouteUseCase(t *testing.T) {
	repo := newFakeRouteRepository()
	repo.route = &domain.Route{
		ID:                       "route-1",
		CompanyID:                "company-1",
		Name:                     "Bogota - Tunja",
		OriginCity:               "Bogota",
		DestinationCity:          "Tunja",
		EstimatedDurationMinutes: 180,
		Status:                   domain.RouteStatusActive,
	}
	uc := NewArchiveRouteUseCase(repo)

	if err := uc.Execute(context.Background(), dto.ArchiveRouteRequest{
		CompanyID: "company-1",
		ID:        "route-1",
	}); err != nil {
		t.Fatalf("archive route: %v", err)
	}

	if repo.archiveCompanyID != "company-1" || repo.archiveID != "route-1" {
		t.Fatalf("repository archive was not company scoped")
	}
}

type fakeRouteRepository struct {
	created *domain.Route
	updated *domain.Route

	route  *domain.Route
	routes []domain.Route
	err    error

	getCompanyID     string
	getID            string
	listCompanyID    string
	archiveCompanyID string
	archiveID        string
}

func newFakeRouteRepository() *fakeRouteRepository {
	return &fakeRouteRepository{}
}

func (r *fakeRouteRepository) Create(_ context.Context, route *domain.Route) error {
	if r.err != nil {
		return r.err
	}
	route.ID = "route-1"
	r.created = route
	return nil
}

func (r *fakeRouteRepository) GetByID(_ context.Context, companyID string, id string) (*domain.Route, error) {
	r.getCompanyID = companyID
	r.getID = id
	if r.err != nil {
		return nil, r.err
	}
	if r.route == nil {
		return nil, domain.ErrRouteNotFound
	}
	return r.route, nil
}

func (r *fakeRouteRepository) ListByCompany(_ context.Context, companyID string) ([]domain.Route, error) {
	r.listCompanyID = companyID
	if r.err != nil {
		return nil, r.err
	}
	return r.routes, nil
}

func (r *fakeRouteRepository) Update(_ context.Context, route *domain.Route) error {
	if r.err != nil {
		return r.err
	}
	r.updated = route
	return nil
}

func (r *fakeRouteRepository) Archive(_ context.Context, companyID string, id string) error {
	r.archiveCompanyID = companyID
	r.archiveID = id
	return r.err
}
