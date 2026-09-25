-- 000001_adm_init_schema.down.sql
-- Rollback Companies & Branches

DROP TABLE IF EXISTS adm_branches CASCADE;
DROP TABLE IF EXISTS adm_companies CASCADE;
