-- 000002_usr_create_users_table.up.sql
-- Enterprise User Management (usr_users)

CREATE TABLE IF NOT EXISTS usr_users (
    id BIGINT PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES adm_companies(id) ON DELETE CASCADE,
    branch_id BIGINT REFERENCES adm_company_branches(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until BIGINT,
    created_at BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at BIGINT,
    updated_by BIGINT,
    CONSTRAINT uq_usr_users_email UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_usr_users_email_lower ON usr_users (LOWER(email));
CREATE INDEX IF NOT EXISTS idx_usr_users_company_id ON usr_users (company_id);
CREATE INDEX IF NOT EXISTS idx_usr_users_branch_id ON usr_users (branch_id);
