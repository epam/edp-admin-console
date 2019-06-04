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
	GetCodebaseByCodebaseAndBranchNames(codebaseName, branchName string) (*models.CodebaseView, error)
}

const (
	SelectCodebaseByNameAndBranchName = "select c.name, cb.name branch_name, c.language, c.build_tool, c.status, al.event branch_status " +
		"from codebase c " +
		"		left join codebase_branch cb on c.id = cb.codebase_id " +
		" 		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on cbal.action_log_id = al.id " +
		"where c.name = ? " +
		"  and cb.name = ? " +
		"  and c.status = 'active'" +
		"  and al.event = 'created';"
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
	codebase := models.CodebaseDetailInfo{Name: codebaseName}

	err := o.Read(&codebase, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "ActionLog")
	codebase.Available = codebase.Status == "active"

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

func (this CodebaseEntityRepository) GetCodebaseByCodebaseAndBranchNames(codebaseName, branchName string) (*models.CodebaseView, error) {
	o := orm.NewOrm()
	var codebase models.CodebaseView

	err := o.Raw(SelectCodebaseByNameAndBranchName, codebaseName, branchName).QueryRow(&codebase)
	if err != nil {
		if err == orm.ErrNoRows {
			log.Printf("Codebase entity wasn't found with %v codebase and %v branch names", codebaseName, branchName)
			return nil, nil
		}
		return nil, err
	}

	return &codebase, nil
}
