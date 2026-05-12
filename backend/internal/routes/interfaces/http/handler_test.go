package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestHandlerCreateUsesAuthenticatedCompany(t *testing.T) {
	create := &fakeCreateRouteUseCase{
		res: &dto.Route{
			ID:                       "route-1",
			CompanyID:                "company-1",
			Name:                     "Bogota - Tunja",
			OriginCity:               "Bogota",
			DestinationCity:          "Tunja",
			EstimatedDistanceKM:      140.5,
			EstimatedDurationMinutes: 180,
			BasePrice:                45000,
			Status:                   string(domain.RouteStatusActive),
		},
	}
	handler := NewHandler(create, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
	req := requestWithClaims(http.MethodPost, "/api/v1/routes", `{
		"company_id": "malicious-company",
		"name": "Bogota - Tunja",
		"origin_city": "Bogota",
		"destination_city": "Tunja",
		"estimated_distance_km": 140.5,
		"estimated_duration_minutes": 180,
		"base_price": 45000,
		"status": "active"
	}`)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if create.req.CompanyID != "company-1" {
		t.Fatalf("expected company from auth context, got %q", create.req.CompanyID)
	}

	var res routeEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Route.ID != "route-1" {
		t.Fatalf("expected route-1, got %q", res.Route.ID)
	}
}

func TestHandlerList(t *testing.T) {
	list := &fakeListRoutesUseCase{
		res: []dto.Route{{
			ID:                       "route-1",
			CompanyID:                "company-1",
			Name:                     "Bogota - Tunja",
			OriginCity:               "Bogota",
			DestinationCity:          "Tunja",
			EstimatedDurationMinutes: 180,
			Status:                   string(domain.RouteStatusActive),
		}},
	}
	handler := NewHandler(&fakeCreateRouteUseCase{}, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, list, &fakeArchiveRouteUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/routes", "")
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if list.req.CompanyID != "company-1" {
		t.Fatalf("expected company from auth context, got %q", list.req.CompanyID)
	}
}

func TestHandlerGetUsesPathIDAndCompanyScope(t *testing.T) {
	get := &fakeGetRouteUseCase{
		res: &dto.Route{
			ID:                       "route-1",
			CompanyID:                "company-1",
			Name:                     "Bogota - Tunja",
			OriginCity:               "Bogota",
			DestinationCity:          "Tunja",
			EstimatedDurationMinutes: 180,
			Status:                   string(domain.RouteStatusActive),
		},
	}
	handler := NewHandler(&fakeCreateRouteUseCase{}, &fakeUpdateRouteUseCase{}, get, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/routes/route-1", "")
	req.SetPathValue("id", "route-1")
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if get.req.CompanyID != "company-1" || get.req.ID != "route-1" {
		t.Fatalf("unexpected get request: %+v", get.req)
	}
}

func TestHandlerUpdate(t *testing.T) {
	update := &fakeUpdateRouteUseCase{
		res: &dto.Route{
			ID:                       "route-1",
			CompanyID:                "company-1",
			Name:                     "Bogota - Duitama",
			OriginCity:               "Bogota",
			DestinationCity:          "Duitama",
			EstimatedDistanceKM:      190.2,
			EstimatedDurationMinutes: 240,
			BasePrice:                60000,
			Status:                   string(domain.RouteStatusInactive),
		},
	}
	handler := NewHandler(&fakeCreateRouteUseCase{}, update, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
	req := requestWithClaims(http.MethodPatch, "/api/v1/routes/route-1", `{
		"name": "Bogota - Duitama",
		"origin_city": "Bogota",
		"destination_city": "Duitama",
		"estimated_distance_km": 190.2,
		"estimated_duration_minutes": 240,
		"base_price": 60000,
		"status": "inactive"
	}`)
	req.SetPathValue("id", "route-1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if update.req.CompanyID != "company-1" || update.req.ID != "route-1" {
		t.Fatalf("unexpected update request: %+v", update.req)
	}
}

func TestHandlerDeleteArchivesRoute(t *testing.T) {
	archive := &fakeArchiveRouteUseCase{}
	handler := NewHandler(&fakeCreateRouteUseCase{}, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, archive)
	req := requestWithClaims(http.MethodDelete, "/api/v1/routes/route-1", "")
	req.SetPathValue("id", "route-1")
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if archive.req.CompanyID != "company-1" || archive.req.ID != "route-1" {
		t.Fatalf("unexpected archive request: %+v", archive.req)
	}
}

func TestHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		body   string
	}{
		{name: "invalid", err: domain.ErrInvalidRoute, status: http.StatusBadRequest, body: "invalid_route"},
		{name: "not found", err: domain.ErrRouteNotFound, status: http.StatusNotFound, body: "route_not_found"},
		{name: "duplicate", err: domain.ErrRouteAlreadyExists, status: http.StatusConflict, body: "route_already_exists"},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, body: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeCreateRouteUseCase{err: tt.err}, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
			req := requestWithClaims(http.MethodPost, "/api/v1/routes", `{}`)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rec.Code)
			}
			assertJSONError(t, rec.Body.String(), tt.body)
		})
	}
}

func TestRoutesRoutesRejectDriver(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateRouteUseCase{}, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{
		claims: &authdomain.TokenClaims{
			UserID:    "driver-1",
			CompanyID: "company-1",
			Role:      userdomain.UserRoleDriver,
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/routes", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "forbidden")
}

func TestRoutesRoutesRequireAuthentication(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateRouteUseCase{}, &fakeUpdateRouteUseCase{}, &fakeGetRouteUseCase{}, &fakeListRoutesUseCase{}, &fakeArchiveRouteUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/routes", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "missing_or_malformed_token")
}

func requestWithClaims(method string, path string, body string) *http.Request {
	var reader *bytes.Buffer
	if body == "" {
		reader = bytes.NewBuffer(nil)
	} else {
		reader = bytes.NewBufferString(body)
	}
	req := httptest.NewRequest(method, path, reader)
	claims := &authdomain.TokenClaims{
		UserID:    "user-1",
		CompanyID: "company-1",
		Role:      userdomain.UserRoleOperator,
	}
	return req.WithContext(authhttp.ContextWithAuthClaims(context.Background(), claims))
}

type fakeCreateRouteUseCase struct {
	req dto.CreateRouteRequest
	res *dto.Route
	err error
}

func (uc *fakeCreateRouteUseCase) Execute(_ context.Context, req dto.CreateRouteRequest) (*dto.Route, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeUpdateRouteUseCase struct {
	req dto.UpdateRouteRequest
	res *dto.Route
	err error
}

func (uc *fakeUpdateRouteUseCase) Execute(_ context.Context, req dto.UpdateRouteRequest) (*dto.Route, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeGetRouteUseCase struct {
	req dto.GetRouteRequest
	res *dto.Route
	err error
}

func (uc *fakeGetRouteUseCase) Execute(_ context.Context, req dto.GetRouteRequest) (*dto.Route, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeListRoutesUseCase struct {
	req dto.ListRoutesRequest
	res []dto.Route
	err error
}

func (uc *fakeListRoutesUseCase) Execute(_ context.Context, req dto.ListRoutesRequest) ([]dto.Route, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeArchiveRouteUseCase struct {
	req dto.ArchiveRouteRequest
	err error
}

func (uc *fakeArchiveRouteUseCase) Execute(_ context.Context, req dto.ArchiveRouteRequest) error {
	uc.req = req
	return uc.err
}

type fakeTokenService struct {
	claims *authdomain.TokenClaims
	err    error
}

func (s *fakeTokenService) Generate(_ context.Context, _ authdomain.TokenSubject) (string, error) {
	return "token", nil
}

func (s *fakeTokenService) Validate(_ context.Context, _ string) (*authdomain.TokenClaims, error) {
	return s.claims, s.err
}

func assertJSONError(t *testing.T, body string, want string) {
	t.Helper()
	if body != `{"error":"`+want+`"}`+"\n" {
		t.Fatalf("expected JSON error %q, got %q", want, body)
	}
}
