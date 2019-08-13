create table applications_to_promote
(
  id serial not null
		constraint applications_to_promote_pk
			primary key,
  cd_pipeline_id integer not null
    constraint cd_pipeline_fk
    references cd_pipeline,
  codebase_id      integer not null
    constraint codebase_fk
    references codebase
);