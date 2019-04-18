package sql_builder

import (
	"edp-admin-console/models"
	"fmt"
)

const (
	SelectAllBranches = "select distinct on (cb.\"name\") cb.name, al.event, al.detailed_message, al.username, al.updated_at " +
		"from codebase_branch cb " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on al.id = cbal.action_log_id " +
		"where c.tenant_name = ? %s " +
		"order by cb.name, al.updated_at desc;"
)

type BranchQueryBuilder struct {
}

func (this *BranchQueryBuilder) GetAllBranchesQuery(filterCriteria models.BranchCriteria) string {
	if filterCriteria.Status == nil {
		return fmt.Sprintf(SelectAllBranches, "")
	}
	if *filterCriteria.Status == "active" {
		return fmt.Sprintf(SelectAllBranches, " and al.event = 'created' ")
	}
	return fmt.Sprintf(SelectAllBranches, " and al.event != 'created' ")
}
