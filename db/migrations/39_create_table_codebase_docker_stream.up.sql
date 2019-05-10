create table if not exists codebase_docker_stream
(
  id                   integer default nextval('codebase_docker_stream_id_seq'::regclass) not null
    constraint codebase_docker_stream_id_pk
      primary key,
  codebase_id          int                                                                not null
    constraint codebase_fk
      references codebase,

  oc_image_stream_name text
);