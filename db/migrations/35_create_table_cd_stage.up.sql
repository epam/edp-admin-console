create table if not exists cd_stage
(
  id                integer default nextval('cd_stage_id_seq'::regclass) not null
    constraint cd_stage_id_pk
      primary key,
  name              text                                                 not null,
  cd_pipeline_id    int                                                  not null
    constraint cd_pipeline_fk
      references cd_pipeline,

  description       text,

  trigger_type      trigger_type                                         not null,
  quality_gate      quality_gate                                         not null,

  jenkins_step_name text,
  "order"           int                                              not null,

  status            status                                               not null,
  constraint cd_stage_pk
    unique (name, cd_pipeline_id)
);