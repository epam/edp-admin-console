/*
 * Copyright 2020 EPAM Systems.
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
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type ICodebaseBranchRepository interface {
	GetCodebaseBranchesByCriteria(criteria query.CodebaseBranchCriteria) ([]query.CodebaseBranch, error)
	SelectDefaultBranchName(appName string) ([]string, error)
}

const (
	SelectDefaultBranchName = "select default_branch from codebase where name = ?;"
)

type CodebaseBranchRepository struct {
	ICodebaseBranchRepository
}

func (CodebaseBranchRepository) GetCodebaseBranchesByCriteria(criteria query.CodebaseBranchCriteria) ([]query.CodebaseBranch, error) {
	o := orm.NewOrm()
	var branches []query.CodebaseBranch

	qs := o.QueryTable(new(query.CodebaseBranch))

	if criteria.Status != "" {
		qs = qs.Filter("status", criteria.Status)
	}

	_, err := qs.OrderBy("name").All(&branches)

	return branches, err
}

func (CodebaseBranchRepository) SelectDefaultBranchName(appName string) ([]string, error) {
	o := orm.NewOrm()
	var defaultBranch []string
	if _, err := o.Raw(SelectDefaultBranchName, appName).QueryRows(&defaultBranch); err != nil {
		return nil, err
	}
	return defaultBranch, nil
}
