-- PostgereSQL 14
-- Database: VS
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