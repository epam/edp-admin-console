create table if not exists jenkins_slave
(
    id   serial not null
        constraint jenkins_slave_pk
            primary key,
    name text   not null unique
);