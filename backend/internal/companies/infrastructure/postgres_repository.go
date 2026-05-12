package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/companies/domain"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return domain.ErrInvalidCompany
	}
	if company.Status == "" {
		company.Status = domain.CompanyStatusActive
	}
	if err := company.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO companies (name, slug, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, company.Name, company.Slug, company.Status).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create company: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	const query = `
		SELECT id, name, slug, status, created_at, updated_at
		FROM companies
		WHERE id = $1`

	return r.get(ctx, query, id)
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	const query = `
		SELECT id, name, slug, status, created_at, updated_at
		FROM companies
		WHERE lower(slug) = lower($1)`

	return r.get(ctx, query, slug)
}

func (r *PostgresRepository) Update(ctx context.Context, company *domain.Company) error {
	if company == nil || company.ID == "" {
		return domain.ErrInvalidCompany
	}
	if err := company.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE companies
		SET name = $2,
			slug = $3,
			status = $4,
			updated_at = now()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, company.ID, company.Name, company.Slug, company.Status).
		Scan(&company.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrCompanyNotFound
	}
	if err != nil {
		return fmt.Errorf("update company: %w", err)
	}

	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.Company, error) {
	var company domain.Company
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&company.ID, &company.Name, &company.Slug, &company.Status, &company.CreatedAt, &company.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompanyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	return &company, nil
}
