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
	"edp-admin-console/models"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego"
	"log"
)

type ApplicationController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    service.BranchService
}

const (
	paramWaitingForBranch = "waitingforbranch"
	paramWaitingForApp    = "waitingforcodebase"
	ApplicationType       = "application"
)

func (this *ApplicationController) GetApplicationsOverviewPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	if flash.Data["success"] != "" {
		this.Data["Success"] = true
	}

	var appType = "application"
	applications, err := this.CodebaseService.GetAllCodebases(models.CodebaseCriteria{
		Type: &appType,
	})
	applications = addCodebaseInProgressIfAny(applications, this.GetString(paramWaitingForApp))
	if err != nil {
		this.Abort("500")
		return
	}

	contextRoles := this.GetSession("realm_roles").([]string)
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(contextRoles)
	this.Data["Codebases"] = applications
	this.Data["Type"] = appType
	this.TplName = "codebase.html"
}

func addCodebaseInProgressIfAny(codebases []models.CodebaseView, codebaseInProgress string) []models.CodebaseView {
	if codebaseInProgress != "" {
		for _, codebase := range codebases {
			if codebase.Name == codebaseInProgress {
				return codebases
			}
		}

		log.Printf("Adding codebase %s which is going to be created to the list.", codebaseInProgress)
		app := models.CodebaseView{
			Name:   codebaseInProgress,
			Status: "in_progress",
		}
		codebases = append(codebases, app)
	}
	return codebases
}

func (this *ApplicationController) GetApplicationOverviewPage() {
	appName := this.GetString(":appName")
	application, err := this.CodebaseService.GetCodebase(appName)
	if err != nil {
		this.Abort("500")
		return
	}

	branchEntities, err := this.BranchService.GetAllReleaseBranchesByAppName(appName)
	branchEntities = addCodebaseBranchInProgressIfAny(branchEntities, this.GetString(paramWaitingForBranch))
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["ReleaseBranches"] = branchEntities
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["Application"] = application
	this.TplName = "codebase_overview.html"
}

func addCodebaseBranchInProgressIfAny(branches []models.ReleaseBranchView, branchInProgress string) []models.ReleaseBranchView {
	if branchInProgress != "" {
		for _, branch := range branches {
			if branch.Name == branchInProgress {
				return branches
			}
		}

		log.Println("Adding branch " + branchInProgress + " which is going to be created to the list.")
		app := models.ReleaseBranchView{
			Name:  branchInProgress,
			Event: "In progress",
		}
		branches = append(branches, app)
	}
	return branches
}

func (this *ApplicationController) GetCreateApplicationPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue()

	if err != nil {
		this.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		this.Data["Error"] = flash.Data["error"]
	}

	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["IsVcsEnabled"] = isVcsEnabled
	this.TplName = "create_application.html"
}

func (this *ApplicationController) CreateApplication() {
	flash := beego.NewFlash()
	codebase := extractApplicationRequestData(this)
	errMsg := validRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/application/create", 302)
		return
	}
	logRequestData(codebase)

	createdObject, err := this.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Application %s name is already exists.", codebase.Name)
			flash.Store(&this.Controller)
			this.Redirect("/admin/edp/application/create", 302)
			return
		}
		this.Abort("500")
		return
	}

	log.Printf("Application object is saved into k8s: %s", createdObject)
	flash.Success("Application object is created.")
	flash.Store(&this.Controller)
	this.Redirect(fmt.Sprintf("/admin/edp/application/overview?%s=%s#codebaseSuccessModal", paramWaitingForApp, codebase.Name), 302)
}

func extractApplicationRequestData(this *ApplicationController) models.Codebase {
	codebase := models.Codebase{
		Lang:      this.GetString("appLang"),
		BuildTool: this.GetString("buildTool"),
		Name:      this.GetString("nameOfApp"),
		Strategy:  this.GetString("strategy"),
		Type:      ApplicationType,
	}

	framework := this.GetString("framework")
	codebase.Framework = &framework

	isMultiModule, _ := this.GetBool("isMultiModule", false)
	if isMultiModule {
		multimoduleApp := fmt.Sprintf("%s-multimodule", *codebase.Framework)
		codebase.Framework = &multimoduleApp
	}

	repoUrl := this.GetString("gitRepoUrl")
	if repoUrl != "" {
		codebase.Repository = &models.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := this.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			codebase.Repository.Login = this.GetString("repoLogin")
			codebase.Repository.Password = this.GetString("repoPassword")
		}
	}

	vcsLogin := this.GetString("vcsLogin")
	vcsPassword := this.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		codebase.Vcs = &models.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}

	needRoute, _ := this.GetBool("needRoute", false)
	if needRoute {
		codebase.Route = &models.Route{
			Site: this.GetString("routeSite"),
		}
		if len(this.GetString("routePath")) > 0 {
			codebase.Route.Path = this.GetString("routePath")
		}
	}

	needDb, _ := this.GetBool("needDb", false)
	if needDb {
		codebase.Database = &models.Database{
			Kind:     this.GetString("database"),
			Version:  this.GetString("dbVersion"),
			Capacity: this.GetString("dbCapacity") + this.GetString("capacityExt"),
			Storage:  this.GetString("dbPersistentStorage"),
		}
	}
	codebase.Username = this.Ctx.Input.Session("username").(string)
	return codebase
}

func isAdmin(contextRoles []string) bool {
	if contextRoles == nil {
		return false
	}
	return util.Contains(contextRoles, beego.AppConfig.String("adminRole"))
}
