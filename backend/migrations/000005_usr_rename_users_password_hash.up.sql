-- 000005_usr_rename_users_password_hash.up.sql
-- Rename password_hash to password on usr_users table

ALTER TABLE usr_users RENAME COLUMN password_hash TO password;
