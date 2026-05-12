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
	"github.com/terraroute/terra-route/backend/internal/users/domain"
)

func TestPostgresRepositoryEnforcesCompanyScope(t *testing.T) {
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
	companyA := createTestCompany(t, ctx, companyRepo, "a")
	companyB := createTestCompany(t, ctx, companyRepo, "b")
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE company_id IN ($1, $2)", companyA.ID, companyB.ID)
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id IN ($1, $2)", companyA.ID, companyB.ID)
	}()

	repo := NewPostgresRepository(db)
	user := &domain.User{
		CompanyID:    companyA.ID,
		Email:        "repo-user-" + safeName(t.Name()) + "@example.com",
		FullName:     "Repo User",
		Role:         domain.UserRoleOperator,
		Status:       domain.UserStatusActive,
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	got, err := repo.GetByID(ctx, companyA.ID, user.ID)
	if err != nil {
		t.Fatalf("get user by id in owning company: %v", err)
	}
	if got.ID != user.ID {
		t.Fatalf("expected user id %q, got %q", user.ID, got.ID)
	}

	_, err = repo.GetByID(ctx, companyB.ID, user.ID)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound for another company scope, got %v", err)
	}

	users, err := repo.ListByCompany(ctx, companyA.ID)
	if err != nil {
		t.Fatalf("list users by company: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user in company scope, got %d", len(users))
	}
}

func createTestCompany(t *testing.T, ctx context.Context, repo *companyinfra.PostgresRepository, suffix string) *companydomain.Company {
	t.Helper()

	company := &companydomain.Company{
		Name:   "Repo Test Company " + suffix,
		Slug:   "repo-test-users-" + safeName(t.Name()) + "-" + suffix,
		Status: companydomain.CompanyStatusActive,
	}
	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company %s: %v", suffix, err)
	}
	return company
}

func safeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
