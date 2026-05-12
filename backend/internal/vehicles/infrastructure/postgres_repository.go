package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/terraroute/terra-route/backend/internal/vehicles/domain"
)

const uniqueViolationCode = "23505"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, vehicle *domain.Vehicle) error {
	if vehicle == nil {
		return domain.ErrInvalidVehicle
	}
	if vehicle.Status == "" {
		vehicle.Status = domain.VehicleStatusActive
	}
	if err := vehicle.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO vehicles (company_id, plate, internal_code, vehicle_type, brand, model, capacity, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		vehicle.CompanyID,
		vehicle.Plate,
		vehicle.InternalCode,
		vehicle.VehicleType,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Capacity,
		string(vehicle.Status),
	).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if isUniqueViolation(err) {
		return domain.ErrVehicleAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, companyID string, id string) (*domain.Vehicle, error) {
	if companyID == "" || id == "" {
		return nil, domain.ErrInvalidVehicle
	}

	const query = `
		SELECT id, company_id, plate, internal_code, vehicle_type, brand, model, capacity, status, created_at, updated_at
		FROM vehicles
		WHERE company_id = $1
			AND id = $2`

	return r.get(ctx, query, companyID, id)
}

func (r *PostgresRepository) ListByCompany(ctx context.Context, companyID string) ([]domain.Vehicle, error) {
	if companyID == "" {
		return nil, domain.ErrInvalidVehicle
	}

	const query = `
		SELECT id, company_id, plate, internal_code, vehicle_type, brand, model, capacity, status, created_at, updated_at
		FROM vehicles
		WHERE company_id = $1
		ORDER BY created_at DESC, id DESC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("list vehicles by company: %w", err)
	}
	defer rows.Close()

	vehicles := make([]domain.Vehicle, 0)
	for rows.Next() {
		var vehicle domain.Vehicle
		if err := scanVehicle(rows, &vehicle); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list vehicles by company rows: %w", err)
	}

	return vehicles, nil
}

func (r *PostgresRepository) Update(ctx context.Context, vehicle *domain.Vehicle) error {
	if vehicle == nil || vehicle.ID == "" {
		return domain.ErrInvalidVehicle
	}
	if err := vehicle.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE vehicles
		SET plate = $3,
			internal_code = $4,
			vehicle_type = $5,
			brand = $6,
			model = $7,
			capacity = $8,
			status = $9,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2
		RETURNING updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		vehicle.CompanyID,
		vehicle.ID,
		vehicle.Plate,
		vehicle.InternalCode,
		vehicle.VehicleType,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Capacity,
		string(vehicle.Status),
	).Scan(&vehicle.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrVehicleNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrVehicleAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}

	return nil
}

func (r *PostgresRepository) MarkInactive(ctx context.Context, companyID string, id string) error {
	if companyID == "" || id == "" {
		return domain.ErrInvalidVehicle
	}

	const query = `
		UPDATE vehicles
		SET status = $3,
			updated_at = now()
		WHERE company_id = $1
			AND id = $2`

	tag, err := r.db.Exec(ctx, query, companyID, id, string(domain.VehicleStatusInactive))
	if err != nil {
		return fmt.Errorf("mark vehicle inactive: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrVehicleNotFound
	}

	return nil
}

func (r *PostgresRepository) get(ctx context.Context, query string, args ...any) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&vehicle.ID,
		&vehicle.CompanyID,
		&vehicle.Plate,
		&vehicle.InternalCode,
		&vehicle.VehicleType,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Capacity,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrVehicleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return &vehicle, nil
}

func scanVehicle(row pgx.Row, vehicle *domain.Vehicle) error {
	if err := row.Scan(
		&vehicle.ID,
		&vehicle.CompanyID,
		&vehicle.Plate,
		&vehicle.InternalCode,
		&vehicle.VehicleType,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Capacity,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan vehicle: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
