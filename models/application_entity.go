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

type BusinessEntity struct {
	Id       int `orm:"pk;column(id)"`
	Name     string
	Tenant   string
	Delition string
	BeType   string

	BeStatus   []*BeStatus     `orm:"reverse(many)"`
	BeProperty []*BeProperties `orm:"reverse(many)"`
}

type BeStatus struct {
	Id             int `orm:"pk;column(id)"`
	Message        string
	LastTimeUpdate string
	UserName       string

	BusinessEntity *BusinessEntity `orm:"rel(fk);column(be_id)"`
	StatusesList   *StatusesList   `orm:"rel(fk);column(status)"`
}

type BeProperties struct {
	Id       int `orm:"pk"`
	Property string
	Value    string

	BusinessEntity *BusinessEntity `orm:"rel(fk);column(be_id)"`
}

type StatusesList struct {
	StatusId   int `orm:"pk"`
	StatusName string
}

func (this *BusinessEntity) TableName() string {
	return "business_entity"
}

func (this *BeStatus) TableName() string {
	return "be_status"
}

func (this *BeProperties) TableName() string {
	return "be_properties"
}

func (this *StatusesList) TableName() string {
	return "statuses_list"
}
