create table cd_pipeline_codebase_branch
(
  cd_pipeline_id   integer not null
    constraint cd_pipeline_fk
      references cd_pipeline,
  codebase_branch_id integer not null
    constraint codebase_branch_fk
      references codebase_branch
);