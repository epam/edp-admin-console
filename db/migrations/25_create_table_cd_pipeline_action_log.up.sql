create table cd_pipeline_action_log
(
  cd_pipeline_id integer not null
    constraint cd_pipeline_fk
      references cd_pipeline,
  action_log_id  integer not null
    constraint action_log_fk
      references action_log
);