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
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

type IApplicationEntityRepository interface {
	GetAllApplications(edpName string) ([]models.Application, error)
	GetApplication(appName string, edpName string) (*models.ApplicationInfo, error)
}

type ApplicationEntityRepository struct {
	IApplicationEntityRepository
}

func (this ApplicationEntityRepository) GetAllApplications(edpName string) ([]models.Application, error) {
	o := orm.NewOrm()
	var applications []models.Application
	var maps []orm.Params

	var query = "select distinct on (\"name\") cb.name, cb.language, cb.build_tool, al.event as status_name " +
		"from codebase cb " +
		"		left join codebase_action_log cal on cb.id = cal.codebase_id " +
		"		left join action_log al on cal.action_log_id = al.id " +
		"where tenant_name = ? " +
		"order by name, al.updated_at desc;"
	_, err := o.Raw(query, edpName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return []models.Application{}, nil
	}

	for _, row := range maps {
		applications = append(applications, models.Application{
			Name:      row["name"].(string),
			Language:  row["language"].(string),
			BuildTool: row["build_tool"].(string),
			Status:    row["status_name"].(string),
		})
	}
	return applications, nil
}

func (this ApplicationEntityRepository) GetApplication(appName string, edpName string) (*models.ApplicationInfo, error) {
	o := orm.NewOrm()
	var application models.ApplicationInfo
	var maps []orm.Params

	var query = "select cb.name, " +
		"       cb.tenant_name       as tenant, " +
		"       cb.type              as be_type, " +
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
		"       al.username          as user_name, " +
		"       al.detailed_message  as message, " +
		"       al.updated_at        as last_time_update " +
		"from codebase cb " +
		"       left join codebase_action_log cal on cb.id = cal.codebase_id " +
		"       left join action_log al on cal.action_log_id = al.id " +
		"where cb.type = 'application' " +
		"  and cb.name = ? " +
		"  and cb.tenant_name = ? " +
		"order by al.updated_at desc limit 1;"

	_, err := o.Raw(query, appName, edpName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for _, row := range maps {
		application = models.ApplicationInfo{
			Name:      row["name"].(string),
			Tenant:    row["tenant"].(string),
			Type:      row["be_type"].(string),
			Status:    row["status_name"].(string),
			Language:  row["language"].(string),
			BuildTool: row["build_tool"].(string),
			Framework: row["framework"].(string),
			Strategy:  row["strategy"].(string),
		}

		if row["user_name"] != nil {
			application.UserName = row["user_name"].(string)
		}
		if row["message"] != nil {
			application.Message = row["message"].(string)
		}

		application.LastTimeUpdate = formatUnixTimestamp(row["last_time_update"].(string))

		application.Available = row["available"] == "active"

		if row["git_url"] != nil {
			application.GitUrl = row["git_url"].(string)
		}

		if row["route_site"] != nil {
			application.RouteSite = row["route_site"].(string)
		}

		if row["route_path"] != nil {
			application.RoutePath = row["route_path"].(string)
		}

		if row["db_kind"] != nil && row["db_version"] != nil && row["db_capacity"] != nil && row["db_storage"] != nil {
			application.DbKind = row["db_kind"].(string)
			application.DbVersion = row["db_version"].(string)
			application.DbCapacity = row["db_capacity"].(string)
			application.DbStorage = row["db_storage"].(string)
		}
	}
	return &application, nil
}

func formatUnixTimestamp(date string) string {
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Println("Couldn't parse dateTime")
		return ""
	}
	return dateTime.String()
}
