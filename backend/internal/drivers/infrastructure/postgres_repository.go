package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/drivers/domain"
)

const uniqueViolationCode = "23505"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, driver *domain.Driver) error {
	if driver == nil {
		return domain.ErrInvalidDriver
	}
	if driver.Status == "" {
		driver.Status = domain.DriverStatusActive
	}
	if err := driver.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO drivers (company_id, user_id, first_name, last_name, document_number, phone, email, license_number, status)
		VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, NULLIF($7, ''), $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		driver.CompanyID,
		driver.UserID,
		driver.FirstName,
		driver.LastName,
		driver.DocumentNumber,
		driver.Phone,
		driver.Email,
		driver.LicenseNumber,
		string(driver.Status),
	).Scan(&driver.ID, &driver.CreatedAt, &driver.UpdatedAt)
	if isUniqueViolation(err) {
		return domain.ErrDriverAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, companyID string, id string) (*domain.Driver, error) {
	if companyID == "" || id == "" {
		return nil, domain.ErrInvalidDriver
	}

	const query = `
		SELECT id, company_id, COALESCE(user_id::text, ''), first_name, last_name, document_number, phone, COALESCE(email, ''), license_number, status, created_at, updated_at
		FROM drivers
		WHERE company_id = $1
			AND id = $2`

	return r.get(ctx, query, companyID, id)
}

func (r *PostgresRepository) ListByCompany(ctx context.Context, companyID string) ([]domain.Driver, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidDriver
	}

	const query = `
		SELECT id, company_id, COALESCE(user_id::text, ''), first_name, last_name, document_number, phone, COALESCE(email, ''), license_number, status, created_at, updated_at
		FROM drivers
		WHERE company_id = $1
		ORDER BY created_at DESC, id DESC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("list drivers by company: %w", err)
	}
	defer rows.Close()

	drivers := make([]domain.Driver, 0)
	for rows.Next() {
		var driver domain.Driver
		if err := scanDriver(rows, &driver); err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list drivers by company rows: %w", err)
	}

	return drivers, nil
}

func (r *PostgresRepository) Update(ctx context.Context, driver *domain.Driver) error {
	if driver == nil || driver.ID == "" {
		return domain.ErrInvalidDriver
	}
	if err := driver.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE drivers
		SET user_id = NULLIF($3, '')::uuid,
			first_name = $4,
			last_name = $5,
			document_number = $6,
			phone = $7,
			email = NULLIF($8, ''),
			license_number = $9,
			status = $10,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2
		RETURNING updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		driver.CompanyID,
		driver.ID,
		driver.UserID,
		driver.FirstName,
		driver.LastName,
		driver.DocumentNumber,
		driver.Phone,
		driver.Email,
		driver.LicenseNumber,
		string(driver.Status),
	).Scan(&driver.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrDriverNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrDriverAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("update driver: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, companyID string, id string) error {
	if companyID == "" || id == "" {
		return domain.ErrInvalidDriver
	}

	const query = `
		UPDATE drivers
		SET status = $3,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2`

	tag, err := r.db.Exec(ctx, query, companyID, id, string(domain.DriverStatusInactive))
	if err != nil {
		return fmt.Errorf("deactivate driver: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrDriverNotFound
	}

	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&driver.ID,
		&driver.CompanyID,
		&driver.UserID,
		&driver.FirstName,
		&driver.LastName,
		&driver.DocumentNumber,
		&driver.Phone,
		&driver.Email,
		&driver.LicenseNumber,
		&driver.Status,
		&driver.CreatedAt,
		&driver.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrDriverNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get driver: %w", err)
	}

	return &driver, nil
}

func scanDriver(row pgx.Row, driver *domain.Driver) error {
	if err := row.Scan(
		&driver.ID,
		&driver.CompanyID,
		&driver.UserID,
		&driver.FirstName,
		&driver.LastName,
		&driver.DocumentNumber,
		&driver.Phone,
		&driver.Email,
		&driver.LicenseNumber,
		&driver.Status,
		&driver.CreatedAt,
		&driver.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan driver: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
