DROP INDEX IF EXISTS users_company_status_idx;
DROP INDEX IF EXISTS users_company_role_idx;
DROP INDEX IF EXISTS users_company_id_idx;
DROP INDEX IF EXISTS users_email_unique_idx;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS companies_status_idx;
DROP INDEX IF EXISTS companies_slug_unique_idx;
DROP TABLE IF EXISTS companies;
