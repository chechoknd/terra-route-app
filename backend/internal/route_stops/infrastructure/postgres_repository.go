package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

const uniqueViolationCode = "23505"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, companyID string, stop *domain.RouteStop) error {
	if companyID == "" || stop == nil {
		return domain.ErrInvalidRouteStop
	}
	if err := stop.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO route_stops (route_id, name, city, stop_order, latitude, longitude)
		SELECT rt.id, $3, $4, $5, $6, $7
		FROM routes rt
		WHERE rt.company_id = $1
			AND rt.id = $2
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, companyID, stop.RouteID, stop.Name, stop.City, stop.StopOrder, stop.Latitude, stop.Longitude).
		Scan(&stop.ID, &stop.CreatedAt, &stop.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrRouteStopNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrRouteStopAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create route stop: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, companyID string, routeID string, id string) (*domain.RouteStop, error) {
	if companyID == "" || routeID == "" || id == "" {
		return nil, domain.ErrInvalidRouteStop
	}

	const query = `
		SELECT rs.id, rs.route_id, rs.name, rs.city, rs.stop_order, rs.latitude, rs.longitude, rs.created_at, rs.updated_at
		FROM route_stops rs
		INNER JOIN routes rt ON rt.id = rs.route_id
		WHERE rt.company_id = $1
			AND rs.route_id = $2
			AND rs.id = $3`

	return r.get(ctx, query, companyID, routeID, id)
}

func (r *PostgresRepository) ListByRoute(ctx context.Context, companyID string, routeID string) ([]domain.RouteStop, error) {
	if companyID == "" || routeID == "" {
		return nil, domain.ErrInvalidRouteStop
	}
	if err := r.ensureRouteAccess(ctx, companyID, routeID); err != nil {
		return nil, err
	}

	const query = `
		SELECT rs.id, rs.route_id, rs.name, rs.city, rs.stop_order, rs.latitude, rs.longitude, rs.created_at, rs.updated_at
		FROM route_stops rs
		INNER JOIN routes rt ON rt.id = rs.route_id
		WHERE rt.company_id = $1
			AND rs.route_id = $2
		ORDER BY rs.stop_order ASC, rs.created_at ASC`

	rows, err := r.db.Query(ctx, query, companyID, routeID)
	if err != nil {
		return nil, fmt.Errorf("list route stops: %w", err)
	}
	defer rows.Close()

	stops := make([]domain.RouteStop, 0)
	for rows.Next() {
		var stop domain.RouteStop
		if err := scanRouteStop(rows, &stop); err != nil {
			return nil, err
		}
		stops = append(stops, stop)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list route stops rows: %w", err)
	}
	return stops, nil
}

func (r *PostgresRepository) ensureRouteAccess(ctx context.Context, companyID string, routeID string) error {
	const query = `
		SELECT 1
		FROM routes
		WHERE company_id = $1
			AND id = $2`

	var exists int
	err := r.db.QueryRow(ctx, query, companyID, routeID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrRouteStopNotFound
	}
	if err != nil {
		return fmt.Errorf("check route access for route stops: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, companyID string, stop *domain.RouteStop) error {
	if companyID == "" || stop == nil || stop.ID == "" {
		return domain.ErrInvalidRouteStop
	}
	if err := stop.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE route_stops rs
		SET name = $4,
			city = $5,
			stop_order = $6,
			latitude = $7,
			longitude = $8,
			updated_at = now()
		FROM routes rt
		WHERE rt.id = rs.route_id
			AND rt.company_id = $1
			AND rs.route_id = $2
			AND rs.id = $3
		RETURNING rs.updated_at`

	err := r.db.QueryRow(ctx, query, companyID, stop.RouteID, stop.ID, stop.Name, stop.City, stop.StopOrder, stop.Latitude, stop.Longitude).
		Scan(&stop.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrRouteStopNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrRouteStopAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("update route stop: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, companyID string, routeID string, id string) error {
	if companyID == "" || routeID == "" || id == "" {
		return domain.ErrInvalidRouteStop
	}

	const query = `
		DELETE FROM route_stops rs
		USING routes rt
		WHERE rt.id = rs.route_id
			AND rt.company_id = $1
			AND rs.route_id = $2
			AND rs.id = $3`

	tag, err := r.db.Exec(ctx, query, companyID, routeID, id)
	if err != nil {
		return fmt.Errorf("delete route stop: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRouteStopNotFound
	}
	return nil
}

func (r *PostgresRepository) Reorder(ctx context.Context, companyID string, routeID string, orderedIDs []string) error {
	if companyID == "" || routeID == "" || len(orderedIDs) == 0 {
		return domain.ErrInvalidRouteStop
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin reorder route stops: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for index, id := range orderedIDs {
		if id == "" {
			return domain.ErrInvalidRouteStop
		}
		const query = `
			UPDATE route_stops rs
			SET stop_order = $4,
				updated_at = now()
			FROM routes rt
			WHERE rt.id = rs.route_id
				AND rt.company_id = $1
				AND rs.route_id = $2
				AND rs.id = $3`
		tag, err := tx.Exec(ctx, query, companyID, routeID, id, 1000000+index)
		if err != nil {
			return fmt.Errorf("reorder route stop: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrRouteStopNotFound
		}
	}

	for index, id := range orderedIDs {
		const query = `
			UPDATE route_stops rs
			SET stop_order = $4,
				updated_at = now()
			FROM routes rt
			WHERE rt.id = rs.route_id
				AND rt.company_id = $1
				AND rs.route_id = $2
				AND rs.id = $3`
		tag, err := tx.Exec(ctx, query, companyID, routeID, id, index+1)
		if err != nil {
			return fmt.Errorf("reorder route stop: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrRouteStopNotFound
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit reorder route stops: %w", err)
	}
	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.RouteStop, error) {
	var stop domain.RouteStop
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stop.ID,
		&stop.RouteID,
		&stop.Name,
		&stop.City,
		&stop.StopOrder,
		&stop.Latitude,
		&stop.Longitude,
		&stop.CreatedAt,
		&stop.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRouteStopNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get route stop: %w", err)
	}
	return &stop, nil
}

func scanRouteStop(row pgx.Row, stop *domain.RouteStop) error {
	if err := row.Scan(
		&stop.ID,
		&stop.RouteID,
		&stop.Name,
		&stop.City,
		&stop.StopOrder,
		&stop.Latitude,
		&stop.Longitude,
		&stop.CreatedAt,
		&stop.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan route stop: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
