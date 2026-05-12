CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT companies_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT companies_slug_not_empty CHECK (length(trim(slug)) > 0),
    CONSTRAINT companies_status_valid CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE UNIQUE INDEX companies_slug_unique_idx ON companies (lower(slug));
CREATE INDEX companies_status_idx ON companies (status);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE RESTRICT,
    email TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_email_not_empty CHECK (length(trim(email)) > 0),
    CONSTRAINT users_full_name_not_empty CHECK (length(trim(full_name)) > 0),
    CONSTRAINT users_password_hash_not_empty CHECK (length(trim(password_hash)) > 0),
    CONSTRAINT users_role_valid CHECK (role IN ('super_admin', 'company_admin', 'operator', 'driver')),
    CONSTRAINT users_status_valid CHECK (status IN ('active', 'inactive', 'suspended')),
    CONSTRAINT users_company_required_for_company_roles CHECK (
        (role = 'super_admin' AND company_id IS NULL)
        OR (role <> 'super_admin' AND company_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email));
CREATE INDEX users_company_id_idx ON users (company_id);
CREATE INDEX users_company_role_idx ON users (company_id, role);
CREATE INDEX users_company_status_idx ON users (company_id, status);
