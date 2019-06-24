create table cd_stage_codebase
(
  cd_stage_id integer not null
    constraint cd_stage_fk
      references cd_stage,
  codebase_id integer not null
    constraint codebase_fk
      references codebase
);