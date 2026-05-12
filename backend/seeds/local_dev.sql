-- TerraRoute local development seed data.
-- LOCAL DEVELOPMENT ONLY. Do not run this against production databases.

WITH demo_company AS (
    INSERT INTO companies (id, name, slug, status)
    VALUES (
        '11111111-1111-4111-8111-111111111111',
        'TerraRoute Demo Company',
        'terraroute-demo',
        'active'
    )
    ON CONFLICT (lower(slug)) DO UPDATE
    SET name = EXCLUDED.name,
        status = EXCLUDED.status,
        updated_at = now()
    RETURNING id
)
INSERT INTO users (id, company_id, email, full_name, role, status, password_hash)
SELECT
    '22222222-2222-4222-8222-222222222222',
    demo_company.id,
    'admin@terraroute.local',
    'Demo Company Admin',
    'company_admin',
    'active',
    '$2a$10$o7iOk9BfgVBVH3aoS8YScejlDcmtxHGAvfHJaJkWE.Axpi0.x2bKK'
FROM demo_company
ON CONFLICT (lower(email)) DO UPDATE
SET company_id = EXCLUDED.company_id,
    full_name = EXCLUDED.full_name,
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    password_hash = EXCLUDED.password_hash,
    updated_at = now();

INSERT INTO users (id, company_id, email, full_name, role, status, password_hash)
SELECT
    '33333333-3333-4333-8333-333333333333',
    companies.id,
    'operator@terraroute.local',
    'Demo Operator',
    'operator',
    'active',
    '$2a$10$iv.oVKaSfX5clADQoRCJueYHAZywthuspSaVbDt/AAK0RQ27NRmnS'
FROM companies
WHERE lower(slug) = lower('terraroute-demo')
ON CONFLICT (lower(email)) DO UPDATE
SET company_id = EXCLUDED.company_id,
    full_name = EXCLUDED.full_name,
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    password_hash = EXCLUDED.password_hash,
    updated_at = now();

INSERT INTO users (id, company_id, email, full_name, role, status, password_hash)
SELECT
    '44444444-4444-4444-8444-444444444444',
    companies.id,
    'driver@terraroute.local',
    'Demo Driver',
    'driver',
    'active',
    '$2a$10$5o8EtnC9zM0XX1nopx.4AeuoGj2gNgIZjZjdJ1k4zWGQ/RtOYz2u6'
FROM companies
WHERE lower(slug) = lower('terraroute-demo')
ON CONFLICT (lower(email)) DO UPDATE
SET company_id = EXCLUDED.company_id,
    full_name = EXCLUDED.full_name,
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    password_hash = EXCLUDED.password_hash,
    updated_at = now();
