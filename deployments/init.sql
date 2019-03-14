SET client_encoding = 'UTF8';
SET client_min_messages = warning;
CREATE TABLE IF NOT EXISTS public.business_entity (
    id SERIAL NOT NULL,
    name character varying(45) NOT NULL,
    tenant character varying(45) NOT NULL,
    delition bigint DEFAULT 0,
    be_type character varying(45) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE (name, tenant, delition, be_type)
);
CREATE TABLE IF NOT EXISTS public.be_properties (
    id SERIAL NOT NULL,
    be_id integer NOT NULL,
    property character varying(45) NOT NULL,
    value character varying(150) NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_ba_properties_to_application FOREIGN KEY (be_id) REFERENCES public.business_entity(id)
);
CREATE TABLE IF NOT EXISTS public.statuses_list (
    status_id SERIAL NOT NULL,
    status_name character varying(45) NOT NULL,

    PRIMARY KEY (status_id)
);
CREATE TABLE IF NOT EXISTS public.be_status (
    id SERIAL NOT NULL,
    be_id integer NOT NULL,
    message text,
    status integer NOT NULL,
    last_time_update bigint,
    user_name character varying(45) NOT NULL,
    available boolean NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_ba_status_to_application FOREIGN KEY (be_id) REFERENCES public.business_entity(id),
    CONSTRAINT fk_ba_status_to_statuses_list FOREIGN KEY (status) REFERENCES public.statuses_list(status_id)
);
INSERT INTO public.statuses_list(status_id, status_name) VALUES (1, 'INITIALIZED'),(2, 'CREATED'),(3, 'IN PROGRESS'),(4, 'FAILED'),(5, 'REOPENED') on conflict (status_id) do nothing;