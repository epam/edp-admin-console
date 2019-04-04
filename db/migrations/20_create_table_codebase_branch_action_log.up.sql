create table codebase_branch_action_log
(
  codebase_branch_id integer not null
    constraint codebase_branch_fk
      references codebase_branch,
  action_log_id      integer not null
    constraint action_log_fk
      references action_log
);