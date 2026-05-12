CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    plate TEXT NOT NULL,
    internal_code TEXT NOT NULL,
    vehicle_type TEXT NOT NULL,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    capacity INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT vehicles_plate_not_empty CHECK (length(trim(plate)) > 0),
    CONSTRAINT vehicles_internal_code_not_empty CHECK (length(trim(internal_code)) > 0),
    CONSTRAINT vehicles_vehicle_type_not_empty CHECK (length(trim(vehicle_type)) > 0),
    CONSTRAINT vehicles_brand_not_empty CHECK (length(trim(brand)) > 0),
    CONSTRAINT vehicles_model_not_empty CHECK (length(trim(model)) > 0),
    CONSTRAINT vehicles_capacity_positive CHECK (capacity > 0),
    CONSTRAINT vehicles_status_valid CHECK (status IN ('active', 'inactive', 'maintenance', 'unavailable'))
);

CREATE UNIQUE INDEX vehicles_company_plate_unique_idx ON vehicles (company_id, lower(plate));
CREATE UNIQUE INDEX vehicles_company_internal_code_unique_idx ON vehicles (company_id, lower(internal_code));
CREATE INDEX vehicles_company_id_idx ON vehicles (company_id);
CREATE INDEX vehicles_company_status_idx ON vehicles (company_id, status);
CREATE INDEX vehicles_company_vehicle_type_idx ON vehicles (company_id, vehicle_type);
