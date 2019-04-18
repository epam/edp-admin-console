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
)

type IReleaseBranchRepository interface {
	GetAllReleaseBranchesByAppName(appName string) ([]models.ReleaseBranchView, error)
	GetAllReleaseBranches(branchFilterCriteria models.BranchCriteria) ([]models.ReleaseBranchView, error)
	GetReleaseBranch(appName, branchName string) (*models.ReleaseBranchView, error)
}

const (
	SelectAllBranches = "select distinct on (cb.\"name\") cb.name, al.event, al.detailed_message, al.username, al.updated_at " +
		"from codebase_branch cb " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on al.id = cbal.action_log_id " +
		"where c.name = ? " +
		"order by cb.name, al.updated_at desc;"
	SelectBranch = "select distinct on (cb.\"name\") cb.name, al.event, al.detailed_message, al.username, al.updated_at " +
		"from codebase_branch cb " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on al.id = cbal.action_log_id " +
		"where c.name = ? " +
		"	and cb.name = ? " +
		"order by cb.name, al.updated_at desc;"
)

type ReleaseBranchRepository struct {
	IReleaseBranchRepository
	QueryManager sql_builder.BranchQueryBuilder
}

func (this ReleaseBranchRepository) GetAllReleaseBranchesByAppName(appName string) ([]models.ReleaseBranchView, error) {
	o := orm.NewOrm()
	var branches []models.ReleaseBranchView

	_, err := o.Raw(SelectAllBranches, appName).QueryRows(&branches)

	if err != nil {
		if err == orm.ErrNoRows {
			log.Printf("No branch entities are found with {%s} appName parameter", appName)
			return nil, nil
		}
		return nil, err
	}

	return branches, nil
}

func (this ReleaseBranchRepository) GetAllReleaseBranches(branchFilterCriteria models.BranchCriteria) ([]models.ReleaseBranchView, error) {
	o := orm.NewOrm()
	var branches []models.ReleaseBranchView

	selectAllBranchesQuery := this.QueryManager.GetAllBranchesQuery(branchFilterCriteria)
	_, err := o.Raw(selectAllBranchesQuery).QueryRows(&branches)

	if err != nil {
		if err == orm.ErrNoRows {
			log.Println("No branch entities are found")
			return nil, nil
		}
		return nil, err
	}

	return branches, nil
}

func (this ReleaseBranchRepository) GetReleaseBranch(appName, branchName string) (*models.ReleaseBranchView, error) {
	o := orm.NewOrm()
	var branch models.ReleaseBranchView

	err := o.Raw(SelectBranch, appName, branchName).QueryRow(&branch)

	if err != nil {
		if err == orm.ErrNoRows {
			log.Printf("No branch entity found with {%s} appName and {%s} branchName parameters", appName, branchName)
			return nil, nil
		}
		return nil, err
	}

	return &branch, nil
}
