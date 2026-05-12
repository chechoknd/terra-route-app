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
	"github.com/terraroute/terra-route/backend/internal/drivers/application/dto"
	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestHandlerCreateUsesAuthenticatedCompany(t *testing.T) {
	create := &fakeCreateDriverUseCase{
		res: &dto.Driver{
			ID:             "driver-1",
			CompanyID:      "company-1",
			UserID:         "user-1",
			FirstName:      "Ana",
			LastName:       "Torres",
			DocumentNumber: "DOC-001",
			Phone:          "+573001112233",
			Email:          "ana@example.test",
			LicenseNumber:  "LIC-001",
			Status:         string(domain.DriverStatusActive),
		},
	}
	handler := NewHandler(create, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
	req := requestWithClaims(http.MethodPost, "/api/v1/drivers", `{
		"company_id": "malicious-company",
		"user_id": "user-1",
		"first_name": "Ana",
		"last_name": "Torres",
		"document_number": "DOC-001",
		"phone": "+573001112233",
		"email": "ana@example.test",
		"license_number": "LIC-001",
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

	var res driverEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Driver.ID != "driver-1" {
		t.Fatalf("expected driver-1, got %q", res.Driver.ID)
	}
}

func TestHandlerList(t *testing.T) {
	list := &fakeListDriversUseCase{
		res: []dto.Driver{{
			ID:             "driver-1",
			CompanyID:      "company-1",
			FirstName:      "Ana",
			LastName:       "Torres",
			DocumentNumber: "DOC-001",
			Phone:          "+573001112233",
			LicenseNumber:  "LIC-001",
			Status:         string(domain.DriverStatusActive),
		}},
	}
	handler := NewHandler(&fakeCreateDriverUseCase{}, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, list, &fakeDeactivateDriverUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/drivers", "")
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
	get := &fakeGetDriverUseCase{
		res: &dto.Driver{
			ID:             "driver-1",
			CompanyID:      "company-1",
			FirstName:      "Ana",
			LastName:       "Torres",
			DocumentNumber: "DOC-001",
			Phone:          "+573001112233",
			LicenseNumber:  "LIC-001",
			Status:         string(domain.DriverStatusActive),
		},
	}
	handler := NewHandler(&fakeCreateDriverUseCase{}, &fakeUpdateDriverUseCase{}, get, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/drivers/driver-1", "")
	req.SetPathValue("id", "driver-1")
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if get.req.CompanyID != "company-1" || get.req.ID != "driver-1" {
		t.Fatalf("unexpected get request: %+v", get.req)
	}
}

func TestHandlerUpdate(t *testing.T) {
	update := &fakeUpdateDriverUseCase{
		res: &dto.Driver{
			ID:             "driver-1",
			CompanyID:      "company-1",
			FirstName:      "Ana Maria",
			LastName:       "Torres",
			DocumentNumber: "DOC-009",
			Phone:          "+573004445566",
			Email:          "ana.maria@example.test",
			LicenseNumber:  "LIC-009",
			Status:         string(domain.DriverStatusSuspended),
		},
	}
	handler := NewHandler(&fakeCreateDriverUseCase{}, update, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
	req := requestWithClaims(http.MethodPatch, "/api/v1/drivers/driver-1", `{
		"first_name": "Ana Maria",
		"last_name": "Torres",
		"document_number": "DOC-009",
		"phone": "+573004445566",
		"email": "ana.maria@example.test",
		"license_number": "LIC-009",
		"status": "suspended"
	}`)
	req.SetPathValue("id", "driver-1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if update.req.CompanyID != "company-1" || update.req.ID != "driver-1" {
		t.Fatalf("unexpected update request: %+v", update.req)
	}
}

func TestHandlerDeleteDeactivatesDriver(t *testing.T) {
	deactivate := &fakeDeactivateDriverUseCase{}
	handler := NewHandler(&fakeCreateDriverUseCase{}, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, deactivate)
	req := requestWithClaims(http.MethodDelete, "/api/v1/drivers/driver-1", "")
	req.SetPathValue("id", "driver-1")
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if deactivate.req.CompanyID != "company-1" || deactivate.req.ID != "driver-1" {
		t.Fatalf("unexpected deactivate request: %+v", deactivate.req)
	}
}

func TestHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		body   string
	}{
		{name: "invalid", err: domain.ErrInvalidDriver, status: http.StatusBadRequest, body: "invalid_driver"},
		{name: "not found", err: domain.ErrDriverNotFound, status: http.StatusNotFound, body: "driver_not_found"},
		{name: "duplicate", err: domain.ErrDriverAlreadyExists, status: http.StatusConflict, body: "driver_already_exists"},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, body: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeCreateDriverUseCase{err: tt.err}, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
			req := requestWithClaims(http.MethodPost, "/api/v1/drivers", `{}`)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rec.Code)
			}
			assertJSONError(t, rec.Body.String(), tt.body)
		})
	}
}

func TestDriversRoutesRejectDriver(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateDriverUseCase{}, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{
		claims: &authdomain.TokenClaims{
			UserID:    "driver-1",
			CompanyID: "company-1",
			Role:      userdomain.UserRoleDriver,
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/drivers", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "forbidden")
}

func TestDriversRoutesRequireAuthentication(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateDriverUseCase{}, &fakeUpdateDriverUseCase{}, &fakeGetDriverUseCase{}, &fakeListDriversUseCase{}, &fakeDeactivateDriverUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/drivers", nil)
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

type fakeCreateDriverUseCase struct {
	req dto.CreateDriverRequest
	res *dto.Driver
	err error
}

func (uc *fakeCreateDriverUseCase) Execute(_ context.Context, req dto.CreateDriverRequest) (*dto.Driver, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeUpdateDriverUseCase struct {
	req dto.UpdateDriverRequest
	res *dto.Driver
	err error
}

func (uc *fakeUpdateDriverUseCase) Execute(_ context.Context, req dto.UpdateDriverRequest) (*dto.Driver, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeGetDriverUseCase struct {
	req dto.GetDriverRequest
	res *dto.Driver
	err error
}

func (uc *fakeGetDriverUseCase) Execute(_ context.Context, req dto.GetDriverRequest) (*dto.Driver, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeListDriversUseCase struct {
	req dto.ListDriversRequest
	res []dto.Driver
	err error
}

func (uc *fakeListDriversUseCase) Execute(_ context.Context, req dto.ListDriversRequest) ([]dto.Driver, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeDeactivateDriverUseCase struct {
	req dto.DeactivateDriverRequest
	err error
}

func (uc *fakeDeactivateDriverUseCase) Execute(_ context.Context, req dto.DeactivateDriverRequest) error {
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
