package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

const uniqueViolationCode = "23505"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, route *domain.Route) error {
	if route == nil {
		return domain.ErrInvalidRoute
	}
	if route.Status == "" {
		route.Status = domain.RouteStatusActive
	}
	if err := route.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO routes (company_id, name, origin_city, destination_city, estimated_distance_km, estimated_duration_minutes, base_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		route.CompanyID,
		route.Name,
		route.OriginCity,
		route.DestinationCity,
		route.EstimatedDistanceKM,
		route.EstimatedDurationMinutes,
		route.BasePrice,
		string(route.Status),
	).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
	if isUniqueViolation(err) {
		return domain.ErrRouteAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create route: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, companyID string, id string) (*domain.Route, error) {
	if companyID == "" || id == "" {
		return nil, domain.ErrInvalidRoute
	}

	const query = `
		SELECT id, company_id, name, origin_city, destination_city, estimated_distance_km::float8, estimated_duration_minutes, base_price::float8, status, created_at, updated_at
		FROM routes
		WHERE company_id = $1
			AND id = $2`

	return r.get(ctx, query, companyID, id)
}

func (r *PostgresRepository) ListByCompany(ctx context.Context, companyID string) ([]domain.Route, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidRoute
	}

	const query = `
		SELECT id, company_id, name, origin_city, destination_city, estimated_distance_km::float8, estimated_duration_minutes, base_price::float8, status, created_at, updated_at
		FROM routes
		WHERE company_id = $1
		ORDER BY created_at DESC, id DESC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("list routes by company: %w", err)
	}
	defer rows.Close()

	routes := make([]domain.Route, 0)
	for rows.Next() {
		var route domain.Route
		if err := scanRoute(rows, &route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list routes by company rows: %w", err)
	}

	return routes, nil
}

func (r *PostgresRepository) Update(ctx context.Context, route *domain.Route) error {
	if route == nil || route.ID == "" {
		return domain.ErrInvalidRoute
	}
	if err := route.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE routes
		SET name = $3,
			origin_city = $4,
			destination_city = $5,
			estimated_distance_km = $6,
			estimated_duration_minutes = $7,
			base_price = $8,
			status = $9,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2
		RETURNING updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		route.CompanyID,
		route.ID,
		route.Name,
		route.OriginCity,
		route.DestinationCity,
		route.EstimatedDistanceKM,
		route.EstimatedDurationMinutes,
		route.BasePrice,
		string(route.Status),
	).Scan(&route.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrRouteNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrRouteAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("update route: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Archive(ctx context.Context, companyID string, id string) error {
	if companyID == "" || id == "" {
		return domain.ErrInvalidRoute
	}

	const query = `
		UPDATE routes
		SET status = $3,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2`

	tag, err := r.db.Exec(ctx, query, companyID, id, string(domain.RouteStatusArchived))
	if err != nil {
		return fmt.Errorf("archive route: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRouteNotFound
	}

	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.Route, error) {
	var route domain.Route
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&route.ID,
		&route.CompanyID,
		&route.Name,
		&route.OriginCity,
		&route.DestinationCity,
		&route.EstimatedDistanceKM,
		&route.EstimatedDurationMinutes,
		&route.BasePrice,
		&route.Status,
		&route.CreatedAt,
		&route.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRouteNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get route: %w", err)
	}

	return &route, nil
}

func scanRoute(row pgx.Row, route *domain.Route) error {
	if err := row.Scan(
		&route.ID,
		&route.CompanyID,
		&route.Name,
		&route.OriginCity,
		&route.DestinationCity,
		&route.EstimatedDistanceKM,
		&route.EstimatedDurationMinutes,
		&route.BasePrice,
		&route.Status,
		&route.CreatedAt,
		&route.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan route: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
