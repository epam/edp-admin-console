insert into applications_to_promote (cd_pipeline_id, codebase_id)
select cpcb.cd_pipeline_id, cb.codebase_id
from cd_pipeline_codebase_branch cpcb
       left join codebase_branch cb on cpcb.codebase_branch_id = cb.id;
