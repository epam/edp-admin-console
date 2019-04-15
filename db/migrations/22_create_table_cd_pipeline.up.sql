create table cd_pipeline
(
  id     integer default nextval('cd_pipeline_id_seq'::regclass) not null
    constraint cd_pipeline_pk
      primary key,
  name   text                                                    not null unique,
  status status                                                  not null
);