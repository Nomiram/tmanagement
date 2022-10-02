-- PostgereSQL 14
-- Database: VS
-- DROP DATABASE IF EXISTS "VS";
CREATE DATABASE "VS"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'Russian_Russia.1251'
    LC_CTYPE = 'Russian_Russia.1251'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

COMMENT ON DATABASE "VS"
    IS 'DB for Vich sys';