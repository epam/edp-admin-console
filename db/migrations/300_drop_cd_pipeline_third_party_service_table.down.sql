create table cd_pipeline_third_party_service
(
    cd_pipeline_id         integer not null
        constraint cd_pipeline_fk
            references cd_pipeline
            on delete cascade,
    third_party_service_id integer not null
        constraint third_party_service_fk
            references third_party_service
);