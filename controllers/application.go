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
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
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
	BranchService    service.CodebaseBranchService
}

const (
	paramWaitingForBranch   = "waitingforbranch"
	paramWaitingForCodebase = "waitingforcodebase"
)

func (c *ApplicationController) GetApplicationsOverviewPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["success"] != "" {
		c.Data["Success"] = true
	}

	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.App,
	})
	applications = addCodebaseInProgressIfAny(applications, c.GetString(paramWaitingForCodebase))
	if err != nil {
		c.Abort("500")
		return
	}

	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(contextRoles)
	c.Data["Codebases"] = applications
	c.Data["Type"] = query.App
	c.TplName = "codebase.html"
}

func addCodebaseInProgressIfAny(codebases []*query.Codebase, codebaseInProgress string) []*query.Codebase {
	if codebaseInProgress != "" {
		for _, codebase := range codebases {
			if codebase.Name == codebaseInProgress {
				return codebases
			}
		}

		log.Printf("Adding codebase %s which is going to be created to the list.", codebaseInProgress)
		app := query.Codebase{
			Name:   codebaseInProgress,
			Status: query.Inactive,
		}
		codebases = append(codebases, &app)
	}
	return codebases
}

func (c *ApplicationController) GetCreateApplicationPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	isVcsEnabled, err := c.EDPTenantService.GetVcsIntegrationValue()

	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.App
	c.TplName = "create_application.html"
}

func (c *ApplicationController) CreateApplication() {
	flash := beego.NewFlash()
	codebase := extractApplicationRequestData(c)
	errMsg := validRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect("/admin/edp/application/create", 302)
		return
	}
	logRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Application %s name is already exists.", codebase.Name)
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/application/create", 302)
			return
		}
		c.Abort("500")
		return
	}

	log.Printf("Application object is saved into k8s: %s", createdObject)
	flash.Success("Application object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/admin/edp/application/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func extractApplicationRequestData(this *ApplicationController) command.CreateCodebase {
	codebase := command.CreateCodebase{
		Lang:      this.GetString("appLang"),
		BuildTool: this.GetString("buildTool"),
		Name:      this.GetString("nameOfApp"),
		Strategy:  this.GetString("strategy"),
		Type:      "application",
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
		codebase.Repository = &command.Repository{
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
		codebase.Vcs = &command.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}

	needRoute, _ := this.GetBool("needRoute", false)
	if needRoute {
		codebase.Route = &command.Route{
			Site: this.GetString("routeSite"),
		}
		if len(this.GetString("routePath")) > 0 {
			codebase.Route.Path = this.GetString("routePath")
		}
	}

	needDb, _ := this.GetBool("needDb", false)
	if needDb {
		codebase.Database = &command.Database{
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
