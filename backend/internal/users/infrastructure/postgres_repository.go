package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/users/domain"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidUser
	}
	if user.Status == "" {
		user.Status = domain.UserStatusActive
	}
	if err := user.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO users (company_id, email, full_name, role, status, password_hash)
		VALUES (NULLIF($1, '')::uuid, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, user.CompanyID, user.Email, user.FullName, user.Role, user.Status, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, companyID string, id string) (*domain.User, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidUser
	}

	const query = `
		SELECT id, company_id, email, full_name, role, status, password_hash, created_at, updated_at
		FROM users
		WHERE company_id = $1
			AND id = $2`

	return r.get(ctx, query, companyID, id)
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, companyID string, email string) (*domain.User, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidUser
	}

	const query = `
		SELECT id, company_id, email, full_name, role, status, password_hash, created_at, updated_at
		FROM users
		WHERE company_id = $1
			AND lower(email) = lower($2)`

	return r.get(ctx, query, companyID, email)
}

func (r *PostgresRepository) ListByCompany(ctx context.Context, companyID string) ([]domain.User, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidUser
	}

	const query = `
		SELECT id, company_id, email, full_name, role, status, password_hash, created_at, updated_at
		FROM users
		WHERE company_id = $1
		ORDER BY created_at DESC, id DESC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("list users by company: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := scanUser(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users by company rows: %w", err)
	}

	return users, nil
}

func (r *PostgresRepository) Update(ctx context.Context, user *domain.User) error {
	if user == nil || user.ID == "" {
		return domain.ErrInvalidUser
	}
	if err := user.Validate(); err != nil {
		return err
	}
	if user.Role == domain.UserRoleSuperAdmin {
		return domain.ErrInvalidUser
	}

	const query = `
		UPDATE users
		SET email = $3,
			full_name = $4,
			role = $5,
			status = $6,
			password_hash = $7,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, user.CompanyID, user.ID, user.Email, user.FullName, user.Role, user.Status, user.PasswordHash).
		Scan(&user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CompanyID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}

func scanUser(row pgx.Row, user *domain.User) error {
	if err := row.Scan(
		&user.ID,
		&user.CompanyID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan user: %w", err)
	}
	return nil
}
