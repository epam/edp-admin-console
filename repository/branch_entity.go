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
	"fmt"
	"github.com/astaxie/beego/orm"
	"log"
)

type IReleaseBranchRepository interface {
	GetAllReleaseBranches(appName, edpName string) ([]models.ReleaseBranch, error)
	GetReleaseBranch(appName, branchName, edpName string) (*models.ReleaseBranch, error)
}

const (
	SelectAllBranches = "select distinct on (cb.\"name\") cb.name, al.event, al.detailed_message, al.username, al.updated_at " +
		"from codebase_branch cb " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on al.id = cbal.action_log_id " +
		"where c.tenant_name = ? " +
		"	and c.name = ? " +
		"order by cb.name, al.updated_at desc;"

	SelectBranch = "select distinct on (cb.\"name\") cb.name, al.event, al.detailed_message, al.username, al.updated_at " +
		"from codebase_branch cb " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on al.id = cbal.action_log_id " +
		"where c.tenant_name = ? " +
		"	and c.name = ? " +
		"	and cb.name = ? " +
		"order by cb.name, al.updated_at desc;"
)

type ReleaseBranchRepository struct {
	IReleaseBranchRepository
}

func (this ReleaseBranchRepository) GetAllReleaseBranches(appName, edpName string) ([]models.ReleaseBranch, error) {
	o := orm.NewOrm()
	var branches []models.ReleaseBranch

	_, err := o.Raw(SelectAllBranches, edpName, appName).QueryRows(&branches)

	if err != nil {
		if err == orm.ErrNoRows {
			log.Printf("No branch entities found with {%s} appName parameter", appName)
			return nil, nil
		}
		return nil, err
	}

	return branches, nil
}

func (this ReleaseBranchRepository) GetReleaseBranch(appName, branchName, edpName string) (*models.ReleaseBranch, error) {
	o := orm.NewOrm()
	var branch models.ReleaseBranch

	err := o.Raw(SelectBranch, edpName, appName, fmt.Sprintf("%s-%s", appName, branchName)).QueryRow(&branch)

	if err != nil {
		if err == orm.ErrNoRows {
			log.Printf("No branch entity found with {%s} appName and {%s} branchName parameters", appName, branchName)
			return nil, nil
		}
		return nil, err
	}

	return &branch, nil
}
