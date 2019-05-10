alter table codebase_branch
  add output_codebase_docker_stream_id int
    constraint output_codebase_docker_stream_fk
      references codebase_docker_stream;