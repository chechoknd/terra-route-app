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
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
	"github.com/terraroute/terra-route/backend/internal/vehicles/application/dto"
	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

func TestHandlerCreateUsesAuthenticatedCompany(t *testing.T) {
	create := &fakeCreateVehicleUseCase{
		res: &dto.Vehicle{
			ID:           "vehicle-1",
			CompanyID:    "company-1",
			Plate:        "ABC123",
			InternalCode: "BUS-001",
			VehicleType:  "bus",
			Brand:        "Mercedes-Benz",
			Model:        "OF-1721",
			Capacity:     42,
			Status:       string(domain.VehicleStatusActive),
		},
	}
	handler := NewHandler(create, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
	req := requestWithClaims(http.MethodPost, "/api/v1/vehicles", `{
		"company_id": "malicious-company",
		"plate": "ABC123",
		"internal_code": "BUS-001",
		"vehicle_type": "bus",
		"brand": "Mercedes-Benz",
		"model": "OF-1721",
		"capacity": 42,
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

	var res vehicleEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Vehicle.ID != "vehicle-1" {
		t.Fatalf("expected vehicle-1, got %q", res.Vehicle.ID)
	}
}

func TestHandlerList(t *testing.T) {
	list := &fakeListVehiclesUseCase{
		res: []dto.Vehicle{{
			ID:           "vehicle-1",
			CompanyID:    "company-1",
			Plate:        "ABC123",
			InternalCode: "BUS-001",
			VehicleType:  "bus",
			Brand:        "Mercedes-Benz",
			Model:        "OF-1721",
			Capacity:     42,
			Status:       string(domain.VehicleStatusActive),
		}},
	}
	handler := NewHandler(&fakeCreateVehicleUseCase{}, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, list, &fakeDeactivateVehicleUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/vehicles", "")
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
	get := &fakeGetVehicleUseCase{
		res: &dto.Vehicle{
			ID:           "vehicle-1",
			CompanyID:    "company-1",
			Plate:        "ABC123",
			InternalCode: "BUS-001",
			VehicleType:  "bus",
			Brand:        "Mercedes-Benz",
			Model:        "OF-1721",
			Capacity:     42,
			Status:       string(domain.VehicleStatusActive),
		},
	}
	handler := NewHandler(&fakeCreateVehicleUseCase{}, &fakeUpdateVehicleUseCase{}, get, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
	req := requestWithClaims(http.MethodGet, "/api/v1/vehicles/vehicle-1", "")
	req.SetPathValue("id", "vehicle-1")
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if get.req.CompanyID != "company-1" || get.req.ID != "vehicle-1" {
		t.Fatalf("unexpected get request: %+v", get.req)
	}
}

func TestHandlerUpdate(t *testing.T) {
	update := &fakeUpdateVehicleUseCase{
		res: &dto.Vehicle{
			ID:           "vehicle-1",
			CompanyID:    "company-1",
			Plate:        "XYZ987",
			InternalCode: "BUS-009",
			VehicleType:  "van",
			Brand:        "Renault",
			Model:        "Master",
			Capacity:     18,
			Status:       string(domain.VehicleStatusMaintenance),
		},
	}
	handler := NewHandler(&fakeCreateVehicleUseCase{}, update, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
	req := requestWithClaims(http.MethodPatch, "/api/v1/vehicles/vehicle-1", `{
		"plate": "XYZ987",
		"internal_code": "BUS-009",
		"vehicle_type": "van",
		"brand": "Renault",
		"model": "Master",
		"capacity": 18,
		"status": "maintenance"
	}`)
	req.SetPathValue("id", "vehicle-1")
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if update.req.CompanyID != "company-1" || update.req.ID != "vehicle-1" {
		t.Fatalf("unexpected update request: %+v", update.req)
	}
}

func TestHandlerDeleteDeactivatesVehicle(t *testing.T) {
	deactivate := &fakeDeactivateVehicleUseCase{}
	handler := NewHandler(&fakeCreateVehicleUseCase{}, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, deactivate)
	req := requestWithClaims(http.MethodDelete, "/api/v1/vehicles/vehicle-1", "")
	req.SetPathValue("id", "vehicle-1")
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if deactivate.req.CompanyID != "company-1" || deactivate.req.ID != "vehicle-1" {
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
		{name: "invalid", err: domain.ErrInvalidVehicle, status: http.StatusBadRequest, body: "invalid_vehicle"},
		{name: "not found", err: domain.ErrVehicleNotFound, status: http.StatusNotFound, body: "vehicle_not_found"},
		{name: "duplicate", err: domain.ErrVehicleAlreadyExists, status: http.StatusConflict, body: "vehicle_already_exists"},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, body: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeCreateVehicleUseCase{err: tt.err}, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
			req := requestWithClaims(http.MethodPost, "/api/v1/vehicles", `{}`)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rec.Code)
			}
			assertJSONError(t, rec.Body.String(), tt.body)
		})
	}
}

func TestVehiclesRoutesRejectDriver(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateVehicleUseCase{}, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{
		claims: &authdomain.TokenClaims{
			UserID:    "driver-1",
			CompanyID: "company-1",
			Role:      userdomain.UserRoleDriver,
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "forbidden")
}

func TestVehiclesRoutesRequireAuthentication(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(&fakeCreateVehicleUseCase{}, &fakeUpdateVehicleUseCase{}, &fakeGetVehicleUseCase{}, &fakeListVehiclesUseCase{}, &fakeDeactivateVehicleUseCase{})
	RegisterRoutes(mux, handler, &fakeTokenService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles", nil)
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

type fakeCreateVehicleUseCase struct {
	req dto.CreateVehicleRequest
	res *dto.Vehicle
	err error
}

func (uc *fakeCreateVehicleUseCase) Execute(_ context.Context, req dto.CreateVehicleRequest) (*dto.Vehicle, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeUpdateVehicleUseCase struct {
	req dto.UpdateVehicleRequest
	res *dto.Vehicle
	err error
}

func (uc *fakeUpdateVehicleUseCase) Execute(_ context.Context, req dto.UpdateVehicleRequest) (*dto.Vehicle, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeGetVehicleUseCase struct {
	req dto.GetVehicleRequest
	res *dto.Vehicle
	err error
}

func (uc *fakeGetVehicleUseCase) Execute(_ context.Context, req dto.GetVehicleRequest) (*dto.Vehicle, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeListVehiclesUseCase struct {
	req dto.ListVehiclesRequest
	res []dto.Vehicle
	err error
}

func (uc *fakeListVehiclesUseCase) Execute(_ context.Context, req dto.ListVehiclesRequest) ([]dto.Vehicle, error) {
	uc.req = req
	return uc.res, uc.err
}

type fakeDeactivateVehicleUseCase struct {
	req dto.DeactivateVehicleRequest
	err error
}

func (uc *fakeDeactivateVehicleUseCase) Execute(_ context.Context, req dto.DeactivateVehicleRequest) error {
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
