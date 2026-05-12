package infrastructure

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/terraroute/terra-route/backend/internal/companies/domain"
	"github.com/terraroute/terra-route/backend/internal/database"
)

func TestPostgresRepositoryCompanyLifecycle(t *testing.T) {
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

	repo := NewPostgresRepository(db)
	company := &domain.Company{
		Name:   "Repo Test Company",
		Slug:   "repo-test-company-" + t.Name(),
		Status: domain.CompanyStatusActive,
	}

	if err := repo.Create(ctx, company); err != nil {
		t.Fatalf("create company: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM companies WHERE id = $1", company.ID)
	}()

	got, err := repo.GetBySlug(ctx, company.Slug)
	if err != nil {
		t.Fatalf("get company by slug: %v", err)
	}
	if got.ID != company.ID {
		t.Fatalf("expected company id %q, got %q", company.ID, got.ID)
	}

	company.Status = domain.CompanyStatusInactive
	if err := repo.Update(ctx, company); err != nil {
		t.Fatalf("update company: %v", err)
	}

	got, err = repo.GetByID(ctx, company.ID)
	if err != nil {
		t.Fatalf("get company by id: %v", err)
	}
	if got.Status != domain.CompanyStatusInactive {
		t.Fatalf("expected status %q, got %q", domain.CompanyStatusInactive, got.Status)
	}

	_, err = repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	if !errors.Is(err, domain.ErrCompanyNotFound) {
		t.Fatalf("expected ErrCompanyNotFound, got %v", err)
	}
}
