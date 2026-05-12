package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestHandlerAddUsesAuthenticatedCompanyAndPathRoute(t *testing.T) {
	add := &fakeAddRouteStopUseCase{res: &dto.RouteStop{ID: "stop-1", RouteID: "route-1", Name: "Terminal Norte", City: "Bogota", StopOrder: 1}}
	handler := NewHandler(add, &fakeUpdateRouteStopUseCase{}, &fakeListRouteStopsUseCase{}, &fakeDeleteRouteStopUseCase{})
	req := requestWithClaims(http.MethodPost, "/api/v1/routes/route-1/stops", `{
		"company_id": "malicious-company",
		"name": "Terminal Norte",
		"city": "Bogota",
		"stop_order": 1,
		"latitude": 4.710989,
		"longitude": -74.072092
	}`)
	req.SetPathValue("id", "route-1")
	rec := httptest.NewRecorder()

	handler.Add(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if add.req.CompanyID != "company-1" || add.req.RouteID != "route-1" {
		t.Fatalf("unexpected add request: %+v", add.req)
	}
}

func TestHandlerListUsesAuthenticatedCompanyAndPathRoute(t *testing.T) {
	list := &fakeListRouteStopsUseCase{res: []dto.RouteStop{{ID: "stop-1", RouteID: "route-1", Name: "Terminal Norte", City: "Bogota", StopOrder: 1}}}
	handler := NewHandler(&fakeAddRouteStopUseCase{}, &fakeUpdateRouteStopUseCase{}, list, &fakeDeleteRouteStopUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/routes/route-1/stops", "")
	req.SetPathValue("id", "route-1")
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if list.req.CompanyID != "company-1" || list.req.RouteID != "route-1" {
		t.Fatalf("unexpected list request: %+v", list.req)
	}
}

func TestHandlerUpdateUsesPathIDs(t *testing.T) {
	update := &fakeUpdateRouteStopUseCase{res: &dto.RouteStop{ID: "stop-1", RouteID: "route-1", Name: "Updated", City: "Tunja", StopOrder: 2}}
	handler := NewHandler(&fakeAddRouteStopUseCase{}, update, &fakeListRouteStopsUseCase{}, &fakeDeleteRouteStopUseCase{})
	req := requestWithClaims(http.MethodPatch, "/api/v1/routes/route-1/stops/stop-1", `{
		"name": "Updated",
		"city": "Tunja",
		"stop_order": 2,
		"latitude": 5.53528,
		"longitude": -73.36778
	}`)
	req.SetPathValue("id", "route-1")
	req.SetPathValue("stopId", "stop-1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if update.req.CompanyID != "company-1" || update.req.RouteID != "route-1" || update.req.ID != "stop-1" {
		t.Fatalf("unexpected update request: %+v", update.req)
	}
}

func TestHandlerDeleteUsesPathIDs(t *testing.T) {
	deleteUC := &fakeDeleteRouteStopUseCase{}
	handler := NewHandler(&fakeAddRouteStopUseCase{}, &fakeUpdateRouteStopUseCase{}, &fakeListRouteStopsUseCase{}, deleteUC)
	req := requestWithClaims(http.MethodDelete, "/api/v1/routes/route-1/stops/stop-1", "")
	req.SetPathValue("id", "route-1")
	req.SetPathValue("stopId", "stop-1")
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if deleteUC.req.CompanyID != "company-1" || deleteUC.req.RouteID != "route-1" || deleteUC.req.ID != "stop-1" {
		t.Fatalf("unexpected delete request: %+v", deleteUC.req)
	}
}

func TestHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		body   string
	}{
		{name: "invalid", err: domain.ErrInvalidRouteStop, status: http.StatusBadRequest, body: "invalid_route_stop"},
		{name: "not found", err: domain.ErrRouteStopNotFound, status: http.StatusNotFound, body: "route_stop_not_found"},
		{name: "duplicate", err: domain.ErrRouteStopAlreadyExists, status: http.StatusConflict, body: "route_stop_already_exists"},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, body: "internal_error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeAddRouteStopUseCase{err: tt.err}, &fakeUpdateRouteStopUseCase{}, &fakeListRouteStopsUseCase{}, &fakeDeleteRouteStopUseCase{})
			req := requestWithClaims(http.MethodPost, "/api/v1/routes/route-1/stops", `{}`)
			req.SetPathValue("id", "route-1")
			rec := httptest.NewRecorder()

			handler.Add(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rec.Code)
			}
			assertJSONError(t, rec.Body.String(), tt.body)
		})
	}
}

func TestRouteStopRoutesRejectDriver(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeAddRouteStopUseCase{}, &fakeUpdateRouteStopUseCase{}, &fakeListRouteStopsUseCase{}, &fakeDeleteRouteStopUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{claims: &authdomain.TokenClaims{UserID: "driver-1", CompanyID: "company-1", Role: userdomain.UserRoleDriver}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/routes/route-1/stops", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "forbidden")
}

func TestRouteStopRoutesRequireAuthentication(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeAddRouteStopUseCase{}, &fakeUpdateRouteStopUseCase{}, &fakeListRouteStopsUseCase{}, &fakeDeleteRouteStopUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/routes/route-1/stops", nil)
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
	claims := &authdomain.TokenClaims{UserID: "user-1", CompanyID: "company-1", Role: userdomain.UserRoleOperator}
	return req.WithContext(authhttp.ContextWithAuthClaims(context.Background(), claims))
}

type fakeAddRouteStopUseCase struct {
	req dto.AddRouteStopRequest
	res *dto.RouteStop
	err error
}

func (uc *fakeAddRouteStopUseCase) Execute(_ context.Context, req dto.AddRouteStopRequest) (*dto.RouteStop, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeUpdateRouteStopUseCase struct {
	req dto.UpdateRouteStopRequest
	res *dto.RouteStop
	err error
}

func (uc *fakeUpdateRouteStopUseCase) Execute(_ context.Context, req dto.UpdateRouteStopRequest) (*dto.RouteStop, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeListRouteStopsUseCase struct {
	req dto.ListRouteStopsRequest
	res []dto.RouteStop
	err error
}

func (uc *fakeListRouteStopsUseCase) Execute(_ context.Context, req dto.ListRouteStopsRequest) ([]dto.RouteStop, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeDeleteRouteStopUseCase struct {
	req dto.DeleteRouteStopRequest
	err error
}

func (uc *fakeDeleteRouteStopUseCase) Execute(_ context.Context, req dto.DeleteRouteStopRequest) error {
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
