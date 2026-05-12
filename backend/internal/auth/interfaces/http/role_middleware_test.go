package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestRequireRolesAllowsConfiguredRole(t *testing.T) {
	req := requestWithRole(userdomain.UserRoleOperator)
	rec := httptest.NewRecorder()

	RequireRoles(userdomain.UserRoleCompanyAdmin, userdomain.UserRoleOperator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestRequireRolesRejectsUnauthenticatedRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	RequireRoles(userdomain.UserRoleOperator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "not_authenticated")
}

func TestRequireRolesRejectsDisallowedRole(t *testing.T) {
	req := requestWithRole(userdomain.UserRoleDriver)
	rec := httptest.NewRecorder()

	RequireRoles(userdomain.UserRoleCompanyAdmin, userdomain.UserRoleOperator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	assertJSONError(t, rec.Body.String(), "forbidden")
}

func TestRequireCompanyAdminOrOperator(t *testing.T) {
	for _, role := range []string{userdomain.UserRoleCompanyAdmin, userdomain.UserRoleOperator} {
		t.Run(role, func(t *testing.T) {
			req := requestWithRole(role)
			rec := httptest.NewRecorder()

			RequireCompanyAdminOrOperator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})).ServeHTTP(rec, req)

			if rec.Code != http.StatusNoContent {
				t.Fatalf("expected status 204, got %d", rec.Code)
			}
		})
	}
}

func TestRequireSuperAdmin(t *testing.T) {
	req := requestWithRole(userdomain.UserRoleSuperAdmin)
	rec := httptest.NewRecorder()

	RequireSuperAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestRequireDriver(t *testing.T) {
	req := requestWithRole(userdomain.UserRoleDriver)
	rec := httptest.NewRecorder()

	RequireDriver(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func requestWithRole(role string) *http.Request {
	claims := &authdomain.TokenClaims{
		UserID:    "user-1",
		CompanyID: "company-1",
		Role:      role,
	}
	return httptest.NewRequest(http.MethodGet, "/protected", nil).
		WithContext(contextWithAuthClaims(context.Background(), claims))
}
