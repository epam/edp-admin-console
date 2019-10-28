create table if not exists edp_component
(
    id   serial not null primary key,
    type text   not null,
    url  text   not null,
    icon text   not null
);