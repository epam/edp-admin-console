package sql_builder

import (
	"edp-admin-console/models"
	"fmt"
)

const (
	SelectAllCodebases = "%s select distinct on (\"name\") cb.name, cb.language, cb.build_tool, al.event as status_name " +
		"from codebase cb " +
		"		left join codebase_action_log cal on cb.id = cal.codebase_id " +
		"		left join action_log al on cal.action_log_id = al.id " +
		"%s " +
		"order by name, al.updated_at desc %s;"
	SelectAllCodebasesWithReleaseBranches = "select c.name as codebase_name, cb.name as branch_name, al.event " +
		"from codebase c " +
		"		left join codebase_branch cb on c.id = cb.codebase_id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on cbal.action_log_id = al.id " +
		"where cb.name is not null %s %s;"
)

func GetAllCodebasesQuery(filterCriteria models.CodebaseCriteria) string {
	var kind string
	if filterCriteria.Type == nil {
		kind = ""
	} else {
		kind = fmt.Sprintf(" where cb.type = '%s' ", *filterCriteria.Type)
	}

	if filterCriteria.Status == nil {
		return fmt.Sprintf(SelectAllCodebases, "", kind, "")
	}
	if *filterCriteria.Status == "active" {
		return fmt.Sprintf(SelectAllCodebases, "select * from(", kind, ") tmp where tmp.status_name = 'created'")
	}
	return fmt.Sprintf(SelectAllCodebases, "select * from(", kind, ") tmp where tmp.status_name != 'created'")
}

func GetAllCodebasesWithReleaseBranchesQuery(filterCriteria models.CodebaseCriteria) string {
	var kind string
	if filterCriteria.Type == nil {
		kind = ""
	} else {
		kind = fmt.Sprintf(" and c.type = '%s' ", *filterCriteria.Type)
	}

	if filterCriteria.Status == nil {
		return fmt.Sprintf(SelectAllCodebasesWithReleaseBranches, "", kind)
	}
	if *filterCriteria.Status == "active" {
		return fmt.Sprintf(SelectAllCodebasesWithReleaseBranches, " and al.event = 'created' ", kind)
	}
	return fmt.Sprintf(SelectAllCodebasesWithReleaseBranches, " and al.event != 'created' ", kind)
}
