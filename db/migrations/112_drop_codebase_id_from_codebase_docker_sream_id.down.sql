alter table codebase_docker_stream
 add codebase_id int not null
  constraint codebase_fk
   references codebase;