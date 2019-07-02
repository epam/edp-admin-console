alter table cd_stage_codebase
 add codebase_branch_id integer not null
  constraint codebase_branch_fk
   references codebase_branch;