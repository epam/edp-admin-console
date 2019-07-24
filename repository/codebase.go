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
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type ICodebaseRepository interface {
	GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error)
	GetCodebaseByName(name string) (*query.Codebase, error)
	GetCodebaseById(id int) (*query.Codebase, error)
	ExistActiveCodebaseAndBranch(codebaseName, branchName string) bool
	ExistCodebaseAndBranch(cbName, brName string) bool
	SelectApplicationToPromote(cdPipelineId int) ([]*query.ApplicationsToPromote, error)
}

type CodebaseRepository struct {
	ICodebaseRepository
}

func (CodebaseRepository) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	o := orm.NewOrm()
	var codebases []*query.Codebase

	qs := o.QueryTable(new(query.Codebase))

	if criteria.Type != "" {
		qs = qs.Filter("type", criteria.Type)
	}
	if criteria.Status != "" {
		qs = qs.Filter("status", criteria.Status)
	}
	_, err := qs.OrderBy("name").
		All(&codebases)

	for _, c := range codebases {

		err = loadRelatedActionLog(c)
		if err != nil {
			return nil, err
		}

		err = loadRelatedCodebaseBranch(c, criteria.BranchStatus)
		if err != nil {
			return nil, err
		}

	}
	return codebases, err
}

func (CodebaseRepository) GetCodebaseByName(name string) (*query.Codebase, error) {
	o := orm.NewOrm()
	codebase := query.Codebase{Name: name}

	err := o.Read(&codebase, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = loadRelatedActionLog(&codebase)

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "CodebaseBranch", false, 100, 0, "Name")

	if err != nil {
		return nil, err
	}

	return &codebase, nil
}

func (CodebaseRepository) GetCodebaseById(id int) (*query.Codebase, error) {
	o := orm.NewOrm()
	codebase := query.Codebase{Id: id}

	err := o.Read(&codebase, "Id")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = loadRelatedActionLog(&codebase)

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "CodebaseBranch", false, 100, 0, "Name")

	if err != nil {
		return nil, err
	}

	return &codebase, nil
}

func loadRelatedActionLog(codebase *query.Codebase) error {
	o := orm.NewOrm()

	_, err := o.QueryTable(new(query.ActionLog)).
		Filter("codebase__codebase_id", codebase.Id).
		OrderBy("LastTimeUpdate").
		Distinct().
		All(&codebase.ActionLog, "LastTimeUpdate", "UserName",
			"Message", "Action", "Result")

	return err
}

func loadRelatedCodebaseBranch(codebase *query.Codebase, status query.Status) error {
	o := orm.NewOrm()

	qs := o.QueryTable(new(query.CodebaseBranch))

	if status != "" {
		qs = qs.Filter("status", status)
	}

	_, err := qs.Filter("codebase_id", codebase.Id).
		OrderBy("Name").
		All(&codebase.CodebaseBranch, "Id", "Name", "FromCommit", "Status")

	return err
}

func (CodebaseRepository) ExistActiveCodebaseAndBranch(cbName, brName string) bool {
	return orm.NewOrm().QueryTable(new(query.Codebase)).
		Filter("name", cbName).
		Filter("status", "active").
		Filter("CodebaseBranch__name", brName).
		Filter("CodebaseBranch__status", "active").
		Exist()
}

func (CodebaseRepository) ExistCodebaseAndBranch(cbName, brName string) bool {
	return orm.NewOrm().QueryTable(new(query.Codebase)).
		Filter("name", cbName).
		Filter("CodebaseBranch__name", brName).
		Exist()
}

func (CodebaseRepository) SelectApplicationToPromote(cdPipelineId int) ([]*query.ApplicationsToPromote, error) {
	o := orm.NewOrm()
	var applicationsToPromote []*query.ApplicationsToPromote

	_, err := o.QueryTable(new(query.ApplicationsToPromote)).
		Filter("cd_pipeline_id", cdPipelineId).All(&applicationsToPromote)

	return applicationsToPromote, err
}
