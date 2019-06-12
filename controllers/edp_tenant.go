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

package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

type EDPTenantController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *EDPTenantController) GetEDPComponents() {
	components := this.EDPTenantService.GetEDPComponents()

	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["InputURL"] = strings.TrimSuffix(this.Ctx.Input.URL(), "/"+context.Tenant)
	this.Data["EDPTenantName"] = context.Tenant
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["EDPComponents"] = components
	this.Data["Type"] = "overview"
	this.TplName = "edp_components.html"
}

func (this *EDPTenantController) GetVcsIntegrationValue() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue()

	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = isVcsEnabled
	this.ServeJSON()
}
