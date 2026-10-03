-- 000003_usr_create_user_sessions_table.up.sql
-- Enterprise User Sessions Table (usr_user_sessions)

CREATE TABLE IF NOT EXISTS usr_user_sessions (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES usr_users(id) ON DELETE CASCADE,
    authorization_token TEXT NOT NULL,
    authorization_token_expired_at BIGINT NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_token_expired_at BIGINT NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT
);

CREATE INDEX IF NOT EXISTS idx_usr_user_sessions_user_id ON usr_user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_usr_user_sessions_refresh_token ON usr_user_sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_usr_user_sessions_auth_token ON usr_user_sessions(authorization_token);
