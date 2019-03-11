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
	"edp-admin-console/models"
	"edp-admin-console/service"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

type ApplicationController struct {
	beego.Controller
	AppService       service.ApplicationService
	EDPTenantService service.EDPTenantService
}

func (this *ApplicationController) GetApplicationsOverviewPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	if flash.Data["success"] != "" {
		this.Data["Success"] = true
	}

	edpTenantName := this.GetString(":name")
	applications, err := this.AppService.GetAllApplications(edpTenantName)
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["CreateApplication"] = fmt.Sprintf("/admin/edp/%s/application/create", this.GetString(":name"))
	this.Data["EdpTenantName"] = edpTenantName
	this.Data["Applications"] = applications
	this.TplName = "application.html"
}

func (this *ApplicationController) GetApplicationOverviewPage() {
	edpTenantName := this.GetString(":name")
	appName := this.GetString(":appName")
	application, err := this.AppService.GetApplication(appName, edpTenantName)
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["Application"] = application
	this.Data["ApplicationsLink"] = fmt.Sprintf("/admin/edp/%s/application/overview", edpTenantName)
	this.TplName = "application_overview.html"
}

func (this *ApplicationController) GetCreateApplicationPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue(this.GetString(":name"))

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if flash.Data["error"] != "" {
		this.Data["Error"] = flash.Data["error"]
	}
	this.Data["IsVcsEnabled"] = isVcsEnabled
	createApplicationLink := fmt.Sprintf("/admin/edp/%s/application", this.GetString(":name"))
	this.Data["CreateApplicationLink"] = createApplicationLink
	this.TplName = "create_application.html"
}

func (this *ApplicationController) CreateApplication() {
	flash := beego.NewFlash()
	edpTenantName := this.GetString(":name")
	app := extractRequestData(this)
	errMsg := validRequestData(app)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect(fmt.Sprintf("/admin/edp/%s/application/create", edpTenantName), 302)
		return
	}

	applicationCr, err := this.AppService.GetApplicationCR(app.Name, edpTenantName)
	if err != nil {
		this.Abort("500")
		return
	}

	application, err := this.AppService.GetApplication(app.Name, edpTenantName)
	if applicationCr != nil || application != nil {
		flash.Error("Application name is already exists.")
		flash.Store(&this.Controller)
		this.Redirect(fmt.Sprintf("/admin/edp/%s/application/create", edpTenantName), 302)
		return
	}

	createdObject, err := this.AppService.CreateApp(app, edpTenantName)

	if err != nil {
		this.Abort("500")
		return
	}

	log.Printf("Application object is saved into k8s: %s", createdObject)
	flash.Success("Application object is created.")
	flash.Store(&this.Controller)
	this.Redirect(fmt.Sprintf("/admin/edp/%s/application/overview", edpTenantName), 302)
}

func extractRequestData(this *ApplicationController) models.App {
	app := models.App{
		Lang:      this.GetString("appLang"),
		Framework: this.GetString("framework"),
		BuildTool: this.GetString("buildTool"),
		Strategy:  this.GetString("strategy"),
		Name:      this.GetString("nameOfApp"),
	}

	isMultiModule, _ := this.GetBool("isMultiModule", false)
	if isMultiModule {
		app.Framework = app.Framework + "-multimodule"
	}

	repoUrl := this.GetString("gitRepoUrl")
	if repoUrl != "" {
		app.Repository = &models.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := this.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			app.Repository.Login = this.GetString("repoLogin")
			app.Repository.Password = this.GetString("repoPassword")
		}
	}

	vcsLogin := this.GetString("vcsLogin")
	vcsPassword := this.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		app.Vcs = &models.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}

	needRoute, _ := this.GetBool("needRoute", false)
	if needRoute {
		app.Route = &models.Route{
			Site: this.GetString("routeSite"),
			Path: this.GetString("routePath"),
		}
	}

	needDb, _ := this.GetBool("needDb", false)
	if needDb {
		app.Database = &models.Database{
			Kind:     this.GetString("database"),
			Version:  this.GetString("dbVersion"),
			Capacity: this.GetString("dbCapacity") + this.GetString("capacityExt"),
			Storage:  this.GetString("dbPersistentStorage"),
		}
	}
	return app
}
