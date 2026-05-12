CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    document_number TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    license_number TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT drivers_first_name_not_empty CHECK (length(trim(first_name)) > 0),
    CONSTRAINT drivers_last_name_not_empty CHECK (length(trim(last_name)) > 0),
    CONSTRAINT drivers_document_number_not_empty CHECK (length(trim(document_number)) > 0),
    CONSTRAINT drivers_phone_not_empty CHECK (length(trim(phone)) > 0),
    CONSTRAINT drivers_email_not_empty_when_present CHECK (email IS NULL OR length(trim(email)) > 0),
    CONSTRAINT drivers_license_number_not_empty CHECK (length(trim(license_number)) > 0),
    CONSTRAINT drivers_status_valid CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE UNIQUE INDEX drivers_company_document_number_unique_idx ON drivers (company_id, lower(document_number));
CREATE UNIQUE INDEX drivers_company_email_unique_idx ON drivers (company_id, lower(email)) WHERE email IS NOT NULL;
CREATE INDEX drivers_company_id_idx ON drivers (company_id);
CREATE INDEX drivers_company_status_idx ON drivers (company_id, status);
CREATE INDEX drivers_company_user_id_idx ON drivers (company_id, user_id) WHERE user_id IS NOT NULL;
CREATE INDEX drivers_company_license_number_idx ON drivers (company_id, lower(license_number));
