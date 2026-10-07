-- 000004_usr_rename_user_sessions_auth_tokens.up.sql
-- Rename authorization_token and authorization_token_expired_at to access_token and access_token_expired_at

ALTER TABLE usr_user_sessions RENAME COLUMN authorization_token TO access_token;
ALTER TABLE usr_user_sessions RENAME COLUMN authorization_token_expired_at TO access_token_expired_at;
