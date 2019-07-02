alter table cd_stage_codebase
 add codebase_id integer not null
  constraint codebase_fk
   references codebase;