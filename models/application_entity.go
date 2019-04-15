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

type Application struct {
	Name      string `json:"name"`
	Language  string `json:"language"`
	BuildTool string `json:"build_tool"`
	Status    string `json:"status"`
}

type ApplicationInfo struct {
	Name           string `json:"name"`
	Tenant         string `json:"tenant"`
	Language       string `json:"language"`
	BuildTool      string `json:"build_tool"`
	Framework      string `json:"framework"`
	Strategy       string `json:"strategy"`
	GitUrl         string `json:"git_url"`
	RouteSite      string `json:"route_site"`
	RoutePath      string `json:"route_path"`
	DbKind         string `json:"db_kind"`
	DbVersion      string `json:"db_version"`
	DbCapacity     string `json:"db_capacity"`
	DbStorage      string `json:"db_storage"`
	DelitionTime   string `json:"delition_time"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	LastTimeUpdate string `json:"last_time_update"`
	UserName       string `json:"user_name"`
	Available      bool   `json:"available"`
	Message        string `json:"message"`
}

type ApplicationWithReleaseBranch struct {
	ApplicationName string
	ReleaseBranches []string
}
