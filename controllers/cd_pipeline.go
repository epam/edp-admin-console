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
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

type CDPipelineController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	PipelineService  service.CDPipelineService
	EDPTenantService service.EDPTenantService
	BranchService    service.BranchService
}

const paramWaitingForCdPipeline = "waitingforcdpipeline"

func (this *CDPipelineController) GetContinuousDeliveryPage() {
	var activeStatus = "active"
	var appType = "application"
	applications, err := this.CodebaseService.GetAllCodebases(models.CodebaseCriteria{
		Status: &activeStatus,
		Type:   &appType,
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
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(contextRoles)
	this.TplName = "continuous_delivery.html"
}

func (this *CDPipelineController) GetCreateCDPipelinePage() {
	flash := beego.ReadFromRequest(&this.Controller)
	var activeStatus = "active"
	var appType = "application"
	applicationsWithReleaseBranches, err := this.CodebaseService.GetAllCodebasesWithReleaseBranches(models.CodebaseCriteria{
		Status: &activeStatus,
		Type:   &appType,
	})
	if err != nil {
		this.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		this.Data["Error"] = flash.Data["error"]
	}

	this.Data["ApplicationsWithReleaseBranches"] = applicationsWithReleaseBranches
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.TplName = "create_cd_pipeline.html"
}

func (this *CDPipelineController) CreateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := this.GetStrings("app")
	pipelineName := this.GetString("pipelineName")
	stages := retrieveStagesFromRequest(this)

	cdPipelineCreateCommand := models.CDPipelineCreateCommand{
		Name:         pipelineName,
		Applications: convertApplicationWithBranchesData(this, appNameCheckboxes),
		Stages:       stages,
	}
	errMsg := validateCDPipelineRequestData(cdPipelineCreateCommand)
	if errMsg != nil {
		log.Printf("Request data is not valid: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/cd-pipeline/create", 302)
		return
	}
	log.Printf("Request data is receieved to create CD pipeline: %s. Applications: %v. Stages: %v",
		cdPipelineCreateCommand.Name, cdPipelineCreateCommand.Applications, cdPipelineCreateCommand.Stages)

	_, pipelineErr := this.PipelineService.CreatePipeline(cdPipelineCreateCommand)
	if pipelineErr != nil {
		if pipelineErr == models.ErrCDPipelineIsExists {
			flash.Error(fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name))
			flash.Store(&this.Controller)
			this.Redirect("/admin/edp/cd-pipeline/create", http.StatusFound)
			return
		}
		if pipelineErr == models.ErrNonValidRelatedBranch {
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications))
			flash.Store(&this.Controller)
			this.Redirect("/admin/edp/cd-pipeline/create", http.StatusBadRequest)
			return
		}
		this.Abort("500")
		return
	}

	this.Data["EDPVersion"] = context.EDPVersion
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

	this.Data["CDPipeline"] = cdPipeline
	this.Data["EDPVersion"] = context.EDPVersion
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

func convertApplicationWithBranchesData(this *CDPipelineController, appNameCheckboxes []string) []models.ApplicationWithBranch {
	var applicationWithBranches []models.ApplicationWithBranch
	for _, appName := range appNameCheckboxes {
		applicationWithBranches = append(applicationWithBranches, models.ApplicationWithBranch{
			ApplicationName: appName,
			BranchName:      this.GetString(appName),
		})
	}
	return applicationWithBranches
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
