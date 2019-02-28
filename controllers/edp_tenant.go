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
	"edp-admin-console/service"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

type EDPTenantController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *EDPTenantController) GetEDPTenants() {
	resourceAccess := this.Ctx.Input.Session("resource_access").(map[string][]string)
	edpTenants, err := this.EDPTenantService.GetEDPTenants(resourceAccess)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["InputURL"] = this.Ctx.Input.URL()
	this.Data["EDPTenants"] = edpTenants
	this.TplName = "edp_tenants.html"
}

func (this *EDPTenantController) GetEDPComponents() {
	resourceAccess := this.Ctx.Input.Session("resource_access").(map[string][]string)
	edpTenants, err := this.EDPTenantService.GetEDPTenants(resourceAccess)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	edpTenantName := this.GetString(":name")
	components := this.EDPTenantService.GetEDPComponents(edpTenantName)
	version, err := this.EDPTenantService.GetEDPVersionByName(edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["InputURL"] = strings.TrimSuffix(this.Ctx.Input.URL(), "/"+edpTenantName)
	appLink := fmt.Sprintf("/admin/edp/%s/application/overview", edpTenantName)
	this.Data["LinkToApplications"] = appLink
	this.Data["EDPTenantName"] = edpTenantName
	this.Data["EDPVersion"] = version
	this.Data["EDPComponents"] = components
	this.Data["EDPTenants"] = edpTenants
	this.TplName = "edp_components.html"
}

func (this *EDPTenantController) GetVcsIntegrationValue() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue(this.GetString(":name"))

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = isVcsEnabled
	this.ServeJSON()
}
