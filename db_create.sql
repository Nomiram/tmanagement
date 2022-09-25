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
	
-- Table: public.orders

-- DROP TABLE IF EXISTS public.orders;

CREATE TABLE IF NOT EXISTS public.orders
(
    order_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    start_date date,
    CONSTRAINT orders_pkey PRIMARY KEY (order_name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;

-- Table: public.tasks

-- DROP TABLE IF EXISTS public.tasks;

CREATE TABLE IF NOT EXISTS public.tasks
(
    task character varying(10) COLLATE pg_catalog."default" NOT NULL,
    order_name character varying(10) COLLATE pg_catalog."default",
    duration integer,
    resource integer,
    pred character varying(10) COLLATE pg_catalog."default",
    CONSTRAINT tasks_pkey PRIMARY KEY (task, order_name),
    CONSTRAINT tasks_order_name_fkey FOREIGN KEY (order_name)
        REFERENCES public.orders (order_name) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;