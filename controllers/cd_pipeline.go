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
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
	"regexp"
)

type CDPipelineController struct {
	beego.Controller
	AppService       service.CodebaseService
	PipelineService  service.CDPipelineService
	EDPTenantService service.EDPTenantService
	BranchService    service.BranchService
}

const paramWaitingForCdPipeline = "waitingforcdpipeline"

func (this *CDPipelineController) GetContinuousDeliveryPage() {
	var activeStatus = "active"
	var appType  = "application"
	applications, err := this.AppService.GetAllCodebases(models.CodebaseCriteria{
		Status: &activeStatus,
		Type: &appType,
	})
	if err != nil {
		this.Abort("500")
		return
	}

	branches, err := this.BranchService.GetAllReleaseBranches(models.BranchCriteria{
		Status: &activeStatus,
	})
	if err != nil {
		this.Abort("500")
		return
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	cdPipelines, err := this.PipelineService.GetAllPipelines(models.CDPipelineCriteria{})
	if err != nil {
		this.Abort("500")
		return
	}
	cdPipelines = addCdPipelineInProgressIfAny(cdPipelines, this.GetString(paramWaitingForCdPipeline))

	contextRoles := this.GetSession("realm_roles").([]string)
	this.Data["ActiveApplicationsAndBranches"] = len(applications) > 0 && len(branches) > 0
	this.Data["CDPipelines"] = cdPipelines
	this.Data["Applications"] = applications
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(contextRoles)
	this.TplName = "continuous_delivery.html"
}

func (this *CDPipelineController) GetCreateCDPipelinePage() {
	flash := beego.ReadFromRequest(&this.Controller)
	var activeStatus = "active"
	var appType = "application"
	applicationsWithReleaseBranches, err := this.AppService.GetAllCodebasesWithReleaseBranches(models.CodebaseCriteria{
		Status: &activeStatus,
		Type:   &appType,
	})
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
	stages := retrieveStagesFromRequest(this)

	errMsg := validateRequestData(appNameCheckboxes, pipelineName, stages)
	if errMsg != nil {
		log.Printf("Request data is not valid: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/cd-pipeline/create", 302)
		return
	}

	releaseBranchCommands := convertRequestReleaseBranchData(appNameCheckboxes, this)
	log.Printf("Request data is receieved to create CD pipelines: %s, with %s name", releaseBranchCommands, pipelineName)

	cdPipeline, err := this.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		this.Abort("500")
		return
	}

	if cdPipeline != nil {
		errMsg := fmt.Sprintf("CD Pipeline %s is already exists.", cdPipeline.Name)
		log.Printf(errMsg)
		flash.Error(errMsg)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/cd-pipeline/create", 302)
		return
	}

	_, err = this.PipelineService.CreatePipeline(pipelineName, releaseBranchCommands)
	if err != nil {
		this.Abort("500")
		return
	}

	_, err = this.PipelineService.CreateStages(pipelineName, stages)
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
	this.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?%s=%s#cdPipelineSuccessModal", paramWaitingForCdPipeline, pipelineName), 302)
}

func (this *CDPipelineController) GetCDPipelineOverviewPage() {
	pipelineName := this.GetString(":pipelineName")

	cdPipeline, err := this.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		this.Abort("500")
		return
	}

	stages, err := this.PipelineService.GetCDPipelineStages(pipelineName)
	if err != nil {
		this.Abort("500")
		return
	}

	cdPipeline.Stages = stages

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["CDPipeline"] = cdPipeline
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.TplName = "cd_pipeline_overview.html"
}

func retrieveStagesFromRequest(this *CDPipelineController) []models.StageCreate {
	var stages []models.StageCreate
	for index, stageName := range this.GetStrings("stageName") {
		stage := models.StageCreate{
			Name:            stageName,
			Description:     this.GetString(stageName + "-stageDesc"),
			StepName:        this.GetString(stageName + "-nameOfStep"),
			QualityGateType: this.GetString(stageName + "-qualityGateType"),
			TriggerType:     this.GetString(stageName + "-triggerType"),
			Order:           index,
		}
		stages = append(stages, stage)
	}
	log.Printf("Stages are fetched from request: %s", stages)
	return stages
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

func validateRequestData(applications []string, pipelineName string, stages []models.StageCreate) *ErrMsg {
	var errorMessage string

	match, _ := regexp.MatchString("^[a-z0-9]([-a-z0-9]*[a-z0-9])$", pipelineName)
	if !match {
		pipelineErrMsg := "Pipeline name may contain only: lower-case letters, numbers, dots and dashes and cannot start and end with dash and dot. Can not be empty."
		log.Println(pipelineErrMsg)
		errorMessage = pipelineErrMsg
	}

	if len(applications) == 0 {
		checkboxErrMsg := "At least one checkbox must be checked"
		log.Println(checkboxErrMsg)
		errorMessage += checkboxErrMsg
	}

	valid := validation.Validation{}

	if len(stages) == 0 {
		stageErrMsg := "At least one stage must be added"
		log.Println(stageErrMsg)
		errorMessage += stageErrMsg
	} else {
		for _, stage := range stages {
			_, err := valid.Valid(stage)
			if err != nil {
				return &ErrMsg{"An internal error has occurred on server while validating stage's form fields.", http.StatusInternalServerError}
			}

			if valid.Errors != nil {
				stageErrMsg := fmt.Sprintf("Stage %s is not valid", stage.Name)
				log.Println(stageErrMsg)
				errorMessage += stageErrMsg
			}
		}
	}

	if len(errorMessage) > 0 {
		return &ErrMsg{
			Message:    errorMessage,
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func addCdPipelineInProgressIfAny(cdPipelines []models.CDPipelineView, pipelineInProgress string) []models.CDPipelineView {
	if pipelineInProgress != "" {
		for _, pipeline := range cdPipelines {
			if pipeline.Name == pipelineInProgress {
				return cdPipelines
			}
		}

		log.Println("Adding CD Pipeline " + pipelineInProgress + " which is going to be created to the list.")

		cdPipelines = append(cdPipelines, models.CDPipelineView{
			Name:   pipelineInProgress,
			Status: "In progress",
		})
	}
	return cdPipelines
}
