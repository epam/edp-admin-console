insert into cd_pipeline_action_log (cd_pipeline_id, action_log_id)
select cs.cd_pipeline_id, csal.action_log_id
from cd_stage_action_log csal
       left join cd_stage cs on csal.cd_stage_id = cs.id;