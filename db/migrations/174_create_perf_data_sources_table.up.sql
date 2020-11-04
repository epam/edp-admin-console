create table if not exists perf_data_sources
(
    id   serial not null
        constraint perf_data_sources_pk
            primary key,
    type text
);