alter table cd_pipeline_third_party_service
  drop constraint cd_pipeline_fk;

alter table cd_pipeline_third_party_service
  add constraint cd_pipeline_fk
    foreign key (cd_pipeline_id) references cd_pipeline
      on delete cascade;

alter table cd_pipeline_action_log
  drop constraint cd_pipeline_fk;

alter table cd_pipeline_action_log
  add constraint cd_pipeline_fk
    foreign key (cd_pipeline_id) references cd_pipeline
      on delete cascade;

alter table cd_stage
  drop constraint cd_pipeline_fk;

alter table cd_stage
  add constraint cd_pipeline_fk
    foreign key (cd_pipeline_id) references cd_pipeline
      on delete cascade;

alter table cd_stage_action_log
  drop constraint cd_stage_fk;

alter table cd_stage_action_log
  add constraint cd_stage_fk
    foreign key (cd_stage_id) references cd_stage
      on delete cascade;

alter table quality_gate_stage
  drop constraint cd_stage_fk;

alter table quality_gate_stage
  add constraint cd_stage_fk
    foreign key (cd_stage_id) references cd_stage
      on delete cascade;

alter table stage_codebase_docker_stream
  drop constraint cd_stage_fk;

alter table stage_codebase_docker_stream
  add constraint cd_stage_fk
    foreign key (cd_stage_id) references cd_stage
      on delete cascade;


alter table stage_codebase_docker_stream
  drop constraint input_codebase_docker_stream_fk;

alter table stage_codebase_docker_stream
  add constraint input_codebase_docker_stream_fk
    foreign key (input_codebase_docker_stream_id) references codebase_docker_stream
      on delete cascade;


alter table stage_codebase_docker_stream
  drop constraint output_codebase_docker_stream_fk;

alter table stage_codebase_docker_stream
  add constraint output_codebase_docker_stream_fk
    foreign key (output_codebase_docker_stream_id) references codebase_docker_stream
      on delete cascade;

alter table applications_to_promote
  drop constraint cd_pipeline_fk;

alter table applications_to_promote
  add constraint cd_pipeline_fk
    foreign key (cd_pipeline_id) references cd_pipeline
      on delete cascade;

alter table cd_pipeline_docker_stream
  drop constraint cd_pipeline_fk;

alter table cd_pipeline_docker_stream
  add constraint cd_pipeline_fk
    foreign key (cd_pipeline_id) references cd_pipeline
      on delete cascade;

alter table cd_pipeline_docker_stream
  drop constraint codebase_docker_stream_fk;

alter table cd_pipeline_docker_stream
  add constraint codebase_docker_stream_fk
    foreign key (codebase_docker_stream_id) references codebase_docker_stream
      on delete cascade;