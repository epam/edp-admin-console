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
)

type MainController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *MainController) Index() {
	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["EDPVersion"] = version
	this.TplName = "index.html"
}

func (this *MainController) CockpitRedirectConfirmationPage() {
	this.Data["Domain"] = fmt.Sprintf("https://edp-cockpit-%s-edp-cicd.%s", this.GetString("edpTenant"), beego.AppConfig.String("dnsWildcard"))
	this.Data["Anchor"] = this.GetString("anchor")
	this.TplName = "cockpit_redirect_confirmation.html"
}
