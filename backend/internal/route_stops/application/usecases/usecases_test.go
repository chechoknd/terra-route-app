package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

func TestAddRouteStopUseCase(t *testing.T) {
	repo := newFakeRouteStopRepository()
	uc := NewAddRouteStopUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.AddRouteStopRequest{
		CompanyID: "company-1",
		RouteID:   "route-1",
		Name:      "Terminal Norte",
		City:      "Bogota",
		StopOrder: 1,
		Latitude:  4.710989,
		Longitude: -74.072092,
	})
	if err != nil {
		t.Fatalf("add route stop: %v", err)
	}
	if res.ID == "" {
		t.Fatal("expected created stop id")
	}
	if repo.createCompanyID != "company-1" || repo.created.RouteID != "route-1" {
		t.Fatalf("repository create was not company scoped")
	}
}

func TestAddRouteStopUseCaseValidatesRequiredFields(t *testing.T) {
	uc := NewAddRouteStopUseCase(newFakeRouteStopRepository())

	_, err := uc.Execute(context.Background(), dto.AddRouteStopRequest{CompanyID: "company-1"})
	if !errors.Is(err, domain.ErrInvalidRouteStop) {
		t.Fatalf("expected ErrInvalidRouteStop, got %v", err)
	}
}

func TestListRouteStopsUseCaseScopesByCompanyAndRoute(t *testing.T) {
	repo := newFakeRouteStopRepository()
	repo.stops = []domain.RouteStop{{ID: "stop-1", RouteID: "route-1", Name: "Terminal Norte", City: "Bogota", StopOrder: 1}}
	uc := NewListRouteStopsUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.ListRouteStopsRequest{CompanyID: "company-1", RouteID: "route-1"})
	if err != nil {
		t.Fatalf("list route stops: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 stop, got %d", len(res))
	}
	if repo.listCompanyID != "company-1" || repo.listRouteID != "route-1" {
		t.Fatalf("repository list was not scoped")
	}
}

func TestUpdateRouteStopUseCase(t *testing.T) {
	repo := newFakeRouteStopRepository()
	repo.stop = &domain.RouteStop{ID: "stop-1", RouteID: "route-1", Name: "Old", City: "Bogota", StopOrder: 1}
	uc := NewUpdateRouteStopUseCase(repo)

	res, err := uc.Execute(context.Background(), dto.UpdateRouteStopRequest{
		CompanyID: "company-1",
		RouteID:   "route-1",
		ID:        "stop-1",
		Name:      "Updated",
		City:      "Tunja",
		StopOrder: 2,
		Latitude:  5.53528,
		Longitude: -73.36778,
	})
	if err != nil {
		t.Fatalf("update route stop: %v", err)
	}
	if res.Name != "Updated" || repo.updated.ID != "stop-1" {
		t.Fatalf("unexpected update result: %+v", res)
	}
}

func TestDeleteRouteStopUseCase(t *testing.T) {
	repo := newFakeRouteStopRepository()
	uc := NewDeleteRouteStopUseCase(repo)

	if err := uc.Execute(context.Background(), dto.DeleteRouteStopRequest{CompanyID: "company-1", RouteID: "route-1", ID: "stop-1"}); err != nil {
		t.Fatalf("delete route stop: %v", err)
	}
	if repo.deleteCompanyID != "company-1" || repo.deleteRouteID != "route-1" || repo.deleteID != "stop-1" {
		t.Fatalf("repository delete was not scoped")
	}
}

func TestReorderRouteStopsUseCase(t *testing.T) {
	repo := newFakeRouteStopRepository()
	uc := NewReorderRouteStopsUseCase(repo)

	if err := uc.Execute(context.Background(), dto.ReorderRouteStopsRequest{CompanyID: "company-1", RouteID: "route-1", OrderedIDs: []string{"stop-2", "stop-1"}}); err != nil {
		t.Fatalf("reorder route stops: %v", err)
	}
	if repo.reorderCompanyID != "company-1" || len(repo.reorderIDs) != 2 {
		t.Fatalf("repository reorder was not scoped")
	}
}

type fakeRouteStopRepository struct {
	created *domain.RouteStop
	updated *domain.RouteStop
	stop    *domain.RouteStop
	stops   []domain.RouteStop
	err     error

	createCompanyID  string
	getCompanyID     string
	getRouteID       string
	getID            string
	listCompanyID    string
	listRouteID      string
	deleteCompanyID  string
	deleteRouteID    string
	deleteID         string
	reorderCompanyID string
	reorderRouteID   string
	reorderIDs       []string
}

func newFakeRouteStopRepository() *fakeRouteStopRepository {
	return &fakeRouteStopRepository{}
}

func (r *fakeRouteStopRepository) Create(_ context.Context, companyID string, stop *domain.RouteStop) error {
	r.createCompanyID = companyID
	if r.err != nil {
		return r.err
	}
	stop.ID = "stop-1"
	r.created = stop
	return nil
}

func (r *fakeRouteStopRepository) GetByID(_ context.Context, companyID string, routeID string, id string) (*domain.RouteStop, error) {
	r.getCompanyID = companyID
	r.getRouteID = routeID
	r.getID = id
	if r.err != nil {
		return nil, r.err
	}
	if r.stop == nil {
		return nil, domain.ErrRouteStopNotFound
	}
	return r.stop, nil
}

func (r *fakeRouteStopRepository) ListByRoute(_ context.Context, companyID string, routeID string) ([]domain.RouteStop, error) {
	r.listCompanyID = companyID
	r.listRouteID = routeID
	return r.stops, r.err
}

func (r *fakeRouteStopRepository) Update(_ context.Context, companyID string, stop *domain.RouteStop) error {
	r.createCompanyID = companyID
	if r.err != nil {
		return r.err
	}
	r.updated = stop
	return nil
}

func (r *fakeRouteStopRepository) Delete(_ context.Context, companyID string, routeID string, id string) error {
	r.deleteCompanyID = companyID
	r.deleteRouteID = routeID
	r.deleteID = id
	return r.err
}

func (r *fakeRouteStopRepository) Reorder(_ context.Context, companyID string, routeID string, orderedIDs []string) error {
	r.reorderCompanyID = companyID
	r.reorderRouteID = routeID
	r.reorderIDs = orderedIDs
	return r.err
}
