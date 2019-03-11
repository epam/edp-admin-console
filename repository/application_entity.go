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
)

type ApplicationEntityRepository struct {
}

func (this ApplicationEntityRepository) GetAllApplications(edpName string) ([]models.BusinessEntity, error) {
	o := orm.NewOrm()
	var applications []models.BusinessEntity

	_, err := o.QueryTable(new(models.BusinessEntity)).Filter("tenant", edpName).All(&applications)
	if err != nil {
		return nil, err
	}

	for i, el := range applications {
		_, err := o.LoadRelated(&el, "BeStatus")
		if err != nil {
			return nil, err
		}
		_, err = o.LoadRelated(&el, "BeProperty")
		if err != nil {
			return nil, err
		}
		for k, status := range el.BeStatus {
			_, err := o.LoadRelated(status, "StatusesList")
			if err != nil {
				return nil, err
			}
			el.BeStatus[k] = status
		}
		applications[i] = el
	}
	return applications, nil
}

func (this ApplicationEntityRepository) GetApplication(appName string, edpName string) (*models.BusinessEntity, error) {
	o := orm.NewOrm()
	var application models.BusinessEntity

	_, err := o.QueryTable(new(models.BusinessEntity)).Filter("name", appName).Filter("tenant", edpName).All(&application)
	if err != nil {
		return nil, err
	}
	_, err = o.LoadRelated(&application, "BeStatus")
	if err != nil {
		return nil, err
	}
	_, err = o.LoadRelated(&application, "BeProperty")
	if err != nil {
		return nil, err
	}
	for k, status := range application.BeStatus {
		_, err := o.LoadRelated(status, "StatusesList")
		if err != nil {
			return nil, err
		}
		application.BeStatus[k] = status
	}

	return &application, nil
}
