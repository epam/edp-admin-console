create table codebase_action_log
(
  codebase_id   integer not null
    constraint codebase_fk
      references codebase,
  action_log_id integer not null
    constraint action_log_fk
      references action_log
);