create table cd_pipeline_docker_stream
(
  cd_pipeline_id   integer not null
    constraint cd_pipeline_fk
      references cd_pipeline,
  codebase_docker_stream_id integer not null
    constraint codebase_docker_stream_fk
      references codebase_docker_stream
);