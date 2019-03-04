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
)

type MainController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *MainController) Index() {
	this.TplName = "index.html"
}

func (this *MainController) GetApplicationPage() {
	createAppLink := fmt.Sprintf("/admin/edp/%s/application/create", this.GetString(":name"))
	this.Data["CreateApplication"] = createAppLink
	this.TplName = "application.html"
}

func (this *MainController) GetCreateApplicationPage() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue(this.GetString(":name"))

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["IsVcsEnabled"] = isVcsEnabled
	this.TplName = "create_application.html"
}
