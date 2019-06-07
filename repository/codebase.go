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
	ExistActiveCodebaseAndBranch(codebaseName, branchName string) bool
	ExistCodebaseAndBranch(cbName, brName string) bool
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
		_, err = o.LoadRelated(c, "ActionLog")
		if err != nil {
			return nil, err
		}

		_, err = o.LoadRelated(c, "CodebaseBranch")
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

	_, err = o.LoadRelated(&codebase, "ActionLog")

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "CodebaseBranch")

	if err != nil {
		return nil, err
	}

	return &codebase, nil
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
