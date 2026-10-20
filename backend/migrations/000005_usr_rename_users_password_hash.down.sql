-- 000005_usr_rename_users_password_hash.down.sql
-- Rollback password to password_hash on usr_users table

ALTER TABLE usr_users RENAME COLUMN password TO password_hash;
