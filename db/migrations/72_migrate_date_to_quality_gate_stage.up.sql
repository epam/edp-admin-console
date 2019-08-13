insert into quality_gate_stage (quality_gate, step_name, cd_stage_id, codebase_id, codebase_branch_id)
select s.quality_gate, s.jenkins_step_name, s.id, cb.codebase_id, cb.id
from cd_stage s
       left join cd_stage_codebase_branch cscb on s.id = cscb.cd_stage_id
       left join codebase_branch cb on cscb.codebase_branch_id = cb.id;