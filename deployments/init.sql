SET client_encoding = 'UTF8';
SET client_min_messages = warning;
CREATE TABLE IF NOT EXISTS public.business_entity (
    name character varying(45) NOT NULL,
    tenant character varying(45) NOT NULL,
    id integer NOT NULL,
    delition timestamp without time zone UNIQUE,
    be_type character varying(45) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE (name, tenant, delition, be_type)
);
CREATE TABLE IF NOT EXISTS public.be_properties (
    be_id integer NOT NULL,
    property character varying(45) NOT NULL,
    value character varying(150) NOT NULL,
    id integer NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_ba_properties_to_application FOREIGN KEY (be_id) REFERENCES public.business_entity(id)
);
CREATE TABLE IF NOT EXISTS public.statuses_list (
    status_id integer NOT NULL,
    status_name character varying(45) NOT NULL,

    PRIMARY KEY (status_id)
);
CREATE TABLE IF NOT EXISTS public.be_status (
    be_id integer NOT NULL,
    message text,
    status integer NOT NULL,
    last_time_update timestamp without time zone,
    user_name character varying(45) NOT NULL,
    id integer NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_ba_status_to_application FOREIGN KEY (be_id) REFERENCES public.business_entity(id),
    CONSTRAINT fk_ba_status_to_statuses_list FOREIGN KEY (status) REFERENCES public.statuses_list(status_id)
);