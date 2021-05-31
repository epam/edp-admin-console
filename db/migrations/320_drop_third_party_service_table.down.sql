create table third_party_service
(
    id          serial not null
        constraint service_pk
            primary key,
    name        text   not null,
    description text   not null,
    version     text   not null,
    url         text,
    icon        text
);