CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    origin_city TEXT NOT NULL,
    destination_city TEXT NOT NULL,
    estimated_distance_km NUMERIC(10,2) NOT NULL,
    estimated_duration_minutes INTEGER NOT NULL,
    base_price NUMERIC(12,2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT routes_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT routes_origin_city_not_empty CHECK (length(trim(origin_city)) > 0),
    CONSTRAINT routes_destination_city_not_empty CHECK (length(trim(destination_city)) > 0),
    CONSTRAINT routes_estimated_distance_non_negative CHECK (estimated_distance_km >= 0),
    CONSTRAINT routes_estimated_duration_positive CHECK (estimated_duration_minutes > 0),
    CONSTRAINT routes_base_price_non_negative CHECK (base_price >= 0),
    CONSTRAINT routes_status_valid CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE UNIQUE INDEX routes_company_name_unique_idx ON routes (company_id, lower(name));
CREATE INDEX routes_company_id_idx ON routes (company_id);
CREATE INDEX routes_company_status_idx ON routes (company_id, status);
CREATE INDEX routes_company_origin_destination_idx ON routes (company_id, lower(origin_city), lower(destination_city));
