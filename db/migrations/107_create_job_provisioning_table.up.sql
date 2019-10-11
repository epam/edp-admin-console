create table if not exists job_provisioning
(
    id   serial not null
        constraint job_provisioning_pk
            primary key,
    name text   not null unique
);