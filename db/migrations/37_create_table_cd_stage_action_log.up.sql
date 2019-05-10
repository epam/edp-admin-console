create table if not exists cd_stage_action_log
(
  cd_stage_id   integer not null
    constraint cd_stage_fk
      references cd_stage,
  action_log_id integer not null
    constraint action_log_fk
      references action_log
);