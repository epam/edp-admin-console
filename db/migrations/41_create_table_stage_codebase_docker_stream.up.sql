create table if not exists stage_codebase_docker_stream
(
  cd_stage_id                      integer not null
    constraint cd_stage_fk
      references cd_stage,
  input_codebase_docker_stream_id  integer not null
    constraint input_codebase_docker_stream_fk
      references codebase_docker_stream,
  output_codebase_docker_stream_id integer
    constraint output_codebase_docker_stream_fk
      references codebase_docker_stream
);