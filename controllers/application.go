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

func (this *ApplicationController) GetCreateApplicationPage() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue(this.GetString(":name"))

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["IsVcsEnabled"] = isVcsEnabled
	createApplicationLink := fmt.Sprintf("/admin/edp/%s/application", this.GetString(":name"))
	this.Data["CreateApplicationLink"] = createApplicationLink
	this.TplName = "create_application.html"
}

func (this *ApplicationController) CreateApplication() {
	app := extractRequestData(this)
	errMsg := validRequestData(app)
	if errMsg != nil {
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}

	edpTenantName := this.GetString(":name")
	createdObject, err := this.AppService.CreateApp(app, edpTenantName)

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, "Failed to create custom resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Custom object is saved into k8s: %s", createdObject)

	appLink := fmt.Sprintf("/admin/edp/%s/application/overview", edpTenantName)
	this.Redirect(appLink, 302)
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
