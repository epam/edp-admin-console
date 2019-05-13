package sql_builder

import (
	"edp-admin-console/models"
	"fmt"
)

const (
	SelectAllCDPipelines = "%s select distinct on (cp.\"name\") cp.name, al.event as status " +
		"from cd_pipeline cp " +
		"		left join cd_pipeline_action_log cpal on cp.id = cpal.cd_pipeline_id " +
		"		left join action_log al on cpal.action_log_id = al.id " +
		"order by cp.name, al.updated_at desc %s;"
)

func GetAllCDPipelinesQuery(filterCriteria models.CDPipelineCriteria) string {
	if filterCriteria.Status == nil {
		return fmt.Sprintf(SelectAllCDPipelines, "", "")
	}
	if *filterCriteria.Status == "active" {
		return fmt.Sprintf(SelectAllCDPipelines, "select * from(", ") tmp where tmp.status = 'created'")
	}
	return fmt.Sprintf(SelectAllCDPipelines, "select * from(", ") tmp where tmp.status != 'created'")
}
