/*
 * Copyright 2019 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository

import (
	"edp-admin-console/models"
	"edp-admin-console/repository/sql_builder"
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

type ICodebaseEntityRepository interface {
	GetAllCodebases(filterCriteria models.CodebaseCriteria) ([]models.CodebaseView, error)
	GetCodebase(codebaseName string) (*models.CodebaseDetailInfo, error)
	GetAllCodebasesWithReleaseBranches(criteria models.CodebaseCriteria) ([]models.CodebaseWithReleaseBranch, error)
}

const (
	SelectCodebase = "select cb.name, " +
		"       cb.type              as cb_type, " +
		"       al.event             as status_name, " +
		"       cb.language, " +
		"       cb.build_tool, " +
		"       cb.framework, " +
		"       cb.strategy, " +
		"       cb.status            as available, " +
		"       cb.repository_url    as git_url, " +
		"       cb.route_site, " +
		"       cb.route_path, " +
		"       cb.database_kind     as db_kind, " +
		"       cb.database_version  as db_version, " +
		"       cb.database_capacity as db_capacity, " +
		"       cb.database_storage  as db_storage, " +
		"       cb.test_report_framework  as test_report_framework, " +
		"       al.username          as user_name, " +
		"       al.detailed_message  as message, " +
		"       al.updated_at        as last_time_update " +
		"from codebase cb " +
		"       left join codebase_action_log cal on cb.id = cal.codebase_id " +
		"       left join action_log al on cal.action_log_id = al.id " +
		"where cb.name = ? " +
		"order by al.updated_at desc limit 1;"
)

type CodebaseEntityRepository struct {
	ICodebaseEntityRepository
}

func (this CodebaseEntityRepository) GetAllCodebases(filterCriteria models.CodebaseCriteria) ([]models.CodebaseView, error) {
	o := orm.NewOrm()
	var codebases []models.CodebaseView
	var maps []orm.Params

	selectAllCodebasesQuery := sql_builder.GetAllCodebasesQuery(filterCriteria)
	_, err := o.Raw(selectAllCodebasesQuery).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return []models.CodebaseView{}, nil
	}

	for _, row := range maps {
		codebases = append(codebases, models.CodebaseView{
			Name:      row["name"].(string),
			Language:  row["language"].(string),
			BuildTool: row["build_tool"].(string),
			Status:    row["status_name"].(string),
		})
	}
	return codebases, nil
}

func (this CodebaseEntityRepository) GetCodebase(codebaseName string) (*models.CodebaseDetailInfo, error) {
	o := orm.NewOrm()
	var codebase models.CodebaseDetailInfo
	var maps []orm.Params

	_, err := o.Raw(SelectCodebase, codebaseName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for _, row := range maps {
		codebase = models.CodebaseDetailInfo{
			Name:      row["name"].(string),
			Type:      row["cb_type"].(string),
			Status:    row["status_name"].(string),
			Language:  row["language"].(string),
			BuildTool: row["build_tool"].(string),
			Framework: row["framework"].(string),
			Strategy:  row["strategy"].(string),
		}

		if row["user_name"] != nil {
			codebase.UserName = row["user_name"].(string)
		}
		if row["message"] != nil {
			codebase.Message = row["message"].(string)
		}

		codebase.LastTimeUpdate = formatUnixTimestamp(row["last_time_update"].(string))

		codebase.Available = row["available"] == "active"

		if row["git_url"] != nil {
			codebase.GitUrl = row["git_url"].(string)
		}

		if row["route_site"] != nil {
			codebase.RouteSite = row["route_site"].(string)
		}

		if row["route_path"] != nil {
			codebase.RoutePath = row["route_path"].(string)
		}

		if row["db_kind"] != nil && row["db_version"] != nil && row["db_capacity"] != nil && row["db_storage"] != nil {
			codebase.DbKind = row["db_kind"].(string)
			codebase.DbVersion = row["db_version"].(string)
			codebase.DbCapacity = row["db_capacity"].(string)
			codebase.DbStorage = row["db_storage"].(string)
		}

		if row["test_report_framework"] != nil {
			codebase.TestReportFramework = row["test_report_framework"].(string)
		}

	}
	return &codebase, nil
}

func (this CodebaseEntityRepository) GetAllCodebasesWithReleaseBranches(criteria models.CodebaseCriteria) ([]models.CodebaseWithReleaseBranch, error) {
	o := orm.NewOrm()
	var codebases []models.CodebaseWithReleaseBranch
	var maps []orm.Params

	query := sql_builder.GetAllCodebasesWithReleaseBranchesQuery(criteria)
	_, err := o.Raw(query).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for _, row := range maps {
		if index, codebase := getCodebase(codebases, row["codebase_name"].(string)); codebase != nil {
			codebase.ReleaseBranches = append(codebase.ReleaseBranches, row["branch_name"].(string))
			codebases[*index] = *codebase
			continue
		}
		codebases = append(codebases, models.CodebaseWithReleaseBranch{
			CodebaseName:    row["codebase_name"].(string),
			ReleaseBranches: []string{row["branch_name"].(string)},
		})
	}

	return codebases, nil
}

func formatUnixTimestamp(date string) string {
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Println("Couldn't parse dateTime")
		return ""
	}
	return dateTime.String()
}

func getCodebase(codebases []models.CodebaseWithReleaseBranch, codebaseName string) (*int, *models.CodebaseWithReleaseBranch) {
	for i, v := range codebases {
		if v.CodebaseName == codebaseName {
			return &i, &v
		}
	}
	return nil, nil
}
