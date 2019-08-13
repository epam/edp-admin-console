create table if not exists quality_gate_stage
(
  id                 integer default nextval('quality_gate_id_seq'::regclass) not null
    constraint quality_gate_pk
      primary key,
  quality_gate       quality_gate                                             not null,
  step_name          text                                                     not null,
  cd_stage_id        int                                                      not null
    constraint cd_stage_fk
      references cd_stage,
  codebase_id        int
    constraint codebase_fk
      references codebase,
  codebase_branch_id int
    constraint codebase_branch_fk
      references codebase_branch
);