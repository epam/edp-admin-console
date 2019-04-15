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
	"github.com/astaxie/beego"
	"log"
	"net/http"
	"regexp"
)

type CDPipelineController struct {
	beego.Controller
	AppService       service.ApplicationService
	PipelineService  service.CDPipelineService
	EDPTenantService service.EDPTenantService
}

func (this *CDPipelineController) GetContinuousDeliveryPage() {
	applications, err := this.AppService.GetAllApplications()
	if err != nil {
		this.Abort("500")
		return
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	contextRoles := this.GetSession("realm_roles").([]string)
	this.Data["Applications"] = applications
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(contextRoles)
	this.TplName = "continuous_delivery.html"
}

func (this *CDPipelineController) GetCreateCDPipelinePage() {
	flash := beego.ReadFromRequest(&this.Controller)
	applicationsWithReleaseBranches, err := this.AppService.GetAllApplicationsWithReleaseBranches()
	if err != nil {
		this.Abort("500")
		return
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		this.Data["Error"] = flash.Data["error"]
	}

	this.Data["ApplicationsWithReleaseBranches"] = applicationsWithReleaseBranches
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.TplName = "create_cd_pipeline.html"
}

func (this *CDPipelineController) CreateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := this.GetStrings("app")
	pipelineName := this.GetString("pipelineName")

	errMsg := validateRequestData(appNameCheckboxes, pipelineName)
	if errMsg != nil {
		log.Printf("Request data is not valid: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/cd-pipeline/create", 302)
		return
	}

	releaseBranchCommands := convertRequestReleaseBranchData(appNameCheckboxes, this)
	log.Printf("Request data is receieved to create CD pipelines: %s, with %s name", releaseBranchCommands, pipelineName)

	_, err := this.PipelineService.CreatePipeline(pipelineName, releaseBranchCommands)
	if err != nil {
		this.Abort("500")
		return
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Redirect("/admin/edp/cd-pipeline/create#cdPipelineSuccessModal", 302)
}

func convertRequestReleaseBranchData(appNameCheckboxes []string, this *CDPipelineController) []models.ReleaseBranchCreatePipelineCommand {
	var releaseBranchCommands []models.ReleaseBranchCreatePipelineCommand
	for _, appName := range appNameCheckboxes {
		releaseBranchCommands = append(releaseBranchCommands, models.ReleaseBranchCreatePipelineCommand{
			AppName:    appName,
			BranchName: this.GetString(appName),
		})
	}
	return releaseBranchCommands
}

func validateRequestData(applications []string, pipelineName string) *ErrMsg {
	var errorMessage string

	match, _ := regexp.MatchString("^[a-z][a-z0-9-.]*[a-z0-9]$", pipelineName)
	if !match {
		pipelineErrMsg := "Pipeline name may contain only: lower-case letters, numbers, dots and dashes and cannot start and end with dash and dot. Can not be empty."
		log.Println(pipelineErrMsg)
		errorMessage = pipelineErrMsg
	}

	if len(applications) == 0 {
		checkboxErrMsg := "At least one checkbox must be checked"
		log.Println(checkboxErrMsg)
		errorMessage = errorMessage + checkboxErrMsg
	}

	if len(errorMessage) > 0 {
		return &ErrMsg{
			Message:    errorMessage,
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}
