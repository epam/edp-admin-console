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
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type IEDPTenantRepository interface {
	GetEdpVersionByName(edpName string) (string, error)
}

type EDPTenantRepository struct {
}

func (this EDPTenantRepository) GetEdpVersionByName(edpName string) (string, error) {
	var edp models.EDPTenant
	o := orm.NewOrm()
	err := o.QueryTable(new(models.EDPTenant)).Filter("name", edpName).One(&edp, "name", "version")
	if err == orm.ErrMultiRows {
		return "", errors.New("the problem has been detected due to db contains more than one value")
	}
	if err == orm.ErrNoRows {
		return "", errors.New(fmt.Sprintf("Couldn't find EDP version for %s.", edpName))
	}
	return edp.Version, nil
}
