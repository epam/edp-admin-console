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

type Codebase struct {
	Name                string      `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	Strategy            string      `json:"strategy"`
	Lang                string      `json:"lang" valid:"Required"`
	Framework           string      `json:"framework" valid:"Required"`
	BuildTool           string      `json:"buildTool" valid:"Required"`
	TestReportFramework *string     `json:"testReportFramework"`
	MultiModule         bool        `json:"multiModule,omitempty"`
	Type                string      `json:"type,omitempty" valid:"Required"`
	Repository          *Repository `json:"repository,omitempty"`
	Route               *Route      `json:"route,omitempty"`
	Database            *Database   `json:"database,omitempty"`
	Vcs                 *Vcs        `json:"vcs,omitempty"`
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
	Name                string `json:"name"`
	Language            string `json:"language"`
	BuildTool           string `json:"build_tool"`
	Framework           string `json:"framework"`
	Strategy            string `json:"strategy"`
	GitUrl              string `json:"git_url"`
	RouteSite           string `json:"route_site"`
	RoutePath           string `json:"route_path"`
	DbKind              string `json:"db_kind"`
	DbVersion           string `json:"db_version"`
	DbCapacity          string `json:"db_capacity"`
	DbStorage           string `json:"db_storage"`
	DelitionTime        string `json:"delition_time"`
	Type                string `json:"type"`
	Status              string `json:"status"`
	LastTimeUpdate      string `json:"last_time_update"`
	UserName            string `json:"user_name"`
	Available           bool   `json:"available"`
	Message             string `json:"message"`
	TestReportFramework string `json:"testReportFramework"`
}

type CodebaseWithReleaseBranch struct {
	CodebaseName    string
	ReleaseBranches []string
}
