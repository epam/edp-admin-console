create table if not exists codebase_perf_data_sources
(
    codebase_id    integer not null
        constraint codebase_fk
            references codebase,
    data_source_id integer not null
        constraint perf_data_sources_fk
            references perf_data_sources
);