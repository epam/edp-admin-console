update codebase_docker_stream
set codebase_branch_id = codebase_branch.id
from codebase_branch
where codebase_branch.output_codebase_docker_stream_id = codebase_docker_stream.id;

update codebase_docker_stream out_cds
set codebase_branch_id = cb.id
from cd_stage cs
       left join stage_codebase_docker_stream scds on cs.id = scds.cd_stage_id
       left join cd_pipeline_codebase_branch cpcb on cs.cd_pipeline_id = cpcb.cd_pipeline_id
       left join codebase_branch cb on cb.id = cpcb.codebase_branch_id
       left join cd_pipeline cp on cs.cd_pipeline_id = cp.id
where out_cds.id = scds.output_codebase_docker_stream_id
  and cb.codebase_id = out_cds.codebase_id;