insert into cd_pipeline_docker_stream (cd_pipeline_id, codebase_docker_stream_id)
select cpcb.cd_pipeline_id, cb.output_codebase_docker_stream_id
from cd_pipeline_codebase_branch cpcb
       left join codebase_branch cb on cpcb.codebase_branch_id = cb.id;