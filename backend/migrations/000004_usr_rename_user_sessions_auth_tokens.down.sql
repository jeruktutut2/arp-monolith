-- 000004_usr_rename_user_sessions_auth_tokens.down.sql
-- Rollback access_token and access_token_expired_at to authorization_token and authorization_token_expired_at

ALTER TABLE usr_user_sessions RENAME COLUMN access_token TO authorization_token;
ALTER TABLE usr_user_sessions RENAME COLUMN access_token_expired_at TO authorization_token_expired_at;
