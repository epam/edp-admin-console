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
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

type CDPipelineController struct {
	beego.Controller
	CodebaseService   service.CodebaseService
	PipelineService   service.CDPipelineService
	EDPTenantService  service.EDPTenantService
	BranchService     service.CodebaseBranchService
	ThirdPartyService service.ThirdPartyService
}

const paramWaitingForCdPipeline = "waitingforcdpipeline"

func (c *CDPipelineController) GetContinuousDeliveryPage() {
	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Status: query.Active,
		Type:   query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	branches, err := c.BranchService.GetCodebaseBranchesByCriteria(query.CodebaseBranchCriteria{
		Status: "active",
	})
	if err != nil {
		c.Abort("500")
		return
	}

	cdPipelines, err := c.PipelineService.GetAllPipelines(query.CDPipelineCriteria{})
	if err != nil {
		c.Abort("500")
		return
	}
	cdPipelines = addCdPipelineInProgressIfAny(cdPipelines, c.GetString(paramWaitingForCdPipeline))

	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["ActiveApplicationsAndBranches"] = len(applications) > 0 && len(branches) > 0
	c.Data["CDPipelines"] = cdPipelines
	c.Data["Applications"] = applications
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(contextRoles)
	c.Data["Type"] = "delivery"
	c.TplName = "continuous_delivery.html"
}

func (c *CDPipelineController) GetCreateCDPipelinePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	apps, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Status: query.Active,
		Type:   query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	services, err := c.ThirdPartyService.GetAllServices()
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["Services"] = services
	c.Data["Apps"] = apps
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Type"] = "delivery"
	c.TplName = "create_cd_pipeline.html"
}

func (c *CDPipelineController) CreateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := c.GetStrings("app")
	pipelineName := c.GetString("pipelineName")
	serviceCheckboxes := c.GetStrings("service")
	stages := retrieveStagesFromRequest(c)

	cdPipelineCreateCommand := models.CDPipelineCreateCommand{
		Name:               pipelineName,
		Applications:       convertApplicationWithBranchesData(c, appNameCheckboxes),
		ThirdPartyServices: serviceCheckboxes,
		Stages:             stages,
	}
	errMsg := validateCDPipelineRequestData(cdPipelineCreateCommand)
	if errMsg != nil {
		log.Printf("Request data is not valid: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect("/admin/edp/cd-pipeline/create", 302)
		return
	}
	log.Printf("Request data is receieved to create CD pipeline: %s. Applications: %v. Stages: %v. Services: %v",
		cdPipelineCreateCommand.Name, cdPipelineCreateCommand.Applications, cdPipelineCreateCommand.Stages, cdPipelineCreateCommand.ThirdPartyServices)

	_, pipelineErr := c.PipelineService.CreatePipeline(cdPipelineCreateCommand)
	if pipelineErr != nil {
		if pipelineErr == models.ErrCDPipelineIsExists {
			flash.Error(fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name))
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/cd-pipeline/create", http.StatusFound)
			return
		}
		if pipelineErr == models.ErrNonValidRelatedBranch {
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications))
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/cd-pipeline/create", http.StatusBadRequest)
			return
		}
		c.Abort("500")
		return
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?%s=%s#cdPipelineSuccessModal", paramWaitingForCdPipeline, pipelineName), 302)
}

func (c *CDPipelineController) GetCDPipelineOverviewPage() {
	pipelineName := c.GetString(":pipelineName")

	cdPipeline, err := c.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["CDPipeline"] = cdPipeline
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Type"] = "delivery"
	c.TplName = "cd_pipeline_overview.html"
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

func addCdPipelineInProgressIfAny(cdPipelines []*query.CDPipeline, pipelineInProgress string) []*query.CDPipeline {
	if pipelineInProgress != "" {
		for _, pipeline := range cdPipelines {
			if pipeline.Name == pipelineInProgress {
				return cdPipelines
			}
		}

		log.Println("Adding CD Pipeline " + pipelineInProgress + " which is going to be created to the list.")
		pipeline := query.CDPipeline{
			Name:   pipelineInProgress,
			Status: "In progress",
		}
		cdPipelines = append(cdPipelines, &pipeline)
	}
	return cdPipelines
}
