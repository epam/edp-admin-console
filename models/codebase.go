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

package models

import (
	"time"
)

type Codebase struct {
	Name                string      `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	Strategy            string      `json:"strategy"`
	Lang                string      `json:"lang" valid:"Required"`
	Framework           *string     `json:"framework,omitempty"`
	BuildTool           string      `json:"buildTool" valid:"Required"`
	TestReportFramework *string     `json:"testReportFramework"`
	MultiModule         bool        `json:"multiModule,omitempty"`
	Type                string      `json:"type,omitempty" valid:"Required"`
	Repository          *Repository `json:"repository,omitempty"`
	Route               *Route      `json:"route,omitempty"`
	Database            *Database   `json:"database,omitempty"`
	Vcs                 *Vcs        `json:"vcs,omitempty"`
	Description         *string     `json:"description,omitempty"`
	Username            string      `json:"username"`
}

type Repository struct {
	Url      string `json:"url,omitempty" valid:"Required;Match(/(?:^git|^ssh|^https?|^git@[-\\w.]+):(\\/\\/)?(.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$/)"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vcs struct {
	Login    string `json:"login,omitempty" valid:"Required"`
	Password string `json:"password,omitempty" valid:"Required"`
}

type Route struct {
	Site string `json:"site,omitempty" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	Path string `json:"path,omitempty" valid:"Match(/^\\/.*$/)"`
}

type Database struct {
	Kind     string `json:"kind,omitempty" valid:"Required"`
	Version  string `json:"version,omitempty" valid:"Required"`
	Capacity string `json:"capacity,omitempty" valid:"Required"`
	Storage  string `json:"storage,omitempty" valid:"Required"`
}

type CodebaseView struct {
	Name      string `json:"name"`
	Language  string `json:"language"`
	BuildTool string `json:"build_tool"`
	Status    string `json:"status"`
}

type CodebaseDetailInfo struct {
	Id                  int          `json:"id" orm:"column(id)"`
	Name                string       `json:"name" orm:"column(name)"`
	Language            string       `json:"language" orm:"column(language)"`
	BuildTool           string       `json:"build_tool" orm:"column(build_tool)"`
	Framework           string       `json:"framework" orm:"column(framework)"`
	Strategy            string       `json:"strategy" orm:"column(strategy)"`
	GitUrl              string       `json:"git_url" orm:"column(repository_url)"`
	RouteSite           string       `json:"route_site" orm:"column(route_site)"`
	RoutePath           string       `json:"route_path" orm:"column(route_path)"`
	DbKind              string       `json:"db_kind" orm:"column(database_kind)"`
	DbVersion           string       `json:"db_version" orm:"column(database_version)"`
	DbCapacity          string       `json:"db_capacity" orm:"column(database_capacity)"`
	DbStorage           string       `json:"db_storage" orm:"column(database_storage)"`
	DelitionTime        string       `json:"delition_time" orm:"-"`
	Type                string       `json:"type" orm:"column(type)"`
	Status              string       `json:"status" orm:"column(status)"`
	Available           bool         `json:"available" orm:"-"`
	TestReportFramework string       `json:"testReportFramework" orm:"column(test_report_framework)"`
	Description         string       `json:"description" orm:"column(description)"`
	ActionLog           []*ActionLog `json:"action_log" orm:"rel(m2m);rel_table(codebase_action_log)"`
}

func (c *CodebaseDetailInfo) TableName() string {
	return "codebase"
}

type ActionLog struct {
	Id             int       `json:"id" orm:"column(id)"`
	LastTimeUpdate time.Time `json:"last_time_update" orm:"column(updated_at)"`
	UserName       string    `json:"user_name" orm:"column(username)"`
	Message        string    `json:"message" orm:"column(detailed_message)"`
	Action         string    `json:"action" orm:"column(action_message)"`
	Result         string    `json:"result" orm:"column(result)"`
}

func (a *ActionLog) TableName() string {
	return "action_log"
}

type CodebaseWithReleaseBranch struct {
	CodebaseName    string
	ReleaseBranches []string
}
