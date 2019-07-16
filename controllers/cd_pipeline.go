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
	"github.com/astaxie/beego/validation"
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
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.App,
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

	autotests, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.Autotests,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["Services"] = services
	c.Data["Apps"] = apps
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Type"] = "delivery"
	c.Data["Autotests"] = autotests
	c.TplName = "create_cd_pipeline.html"
}

func (c *CDPipelineController) GetEditCDPipelinePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	pipelineName := c.GetString(":name")

	cdPipeline, err := c.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		c.Abort("500")
		return
	}

	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}
	c.Data["CDPipeline"] = cdPipeline
	c.Data["Apps"] = applications
	c.Data["Type"] = "delivery"
	c.TplName = "edit_cd_pipeline.html"
}

func (c *CDPipelineController) UpdateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := c.GetStrings("app")
	pipelineName := c.GetString(":name")

	pipelineUpdateCommand := models.CDPipelineCommand{
		Name:         pipelineName,
		Applications: convertApplicationWithBranchesData(c, appNameCheckboxes),
	}

	errMsg := validateCDPipelineUpdateRequestData(pipelineUpdateCommand)
	if errMsg != nil {
		log.Printf("Request data is not valid: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/%s/update", pipelineName), 302)
		return
	}
	log.Printf("Request data is receieved to update CD pipeline: %s. Applications: %v. Stages: %v. Services: %v",
		pipelineName, pipelineUpdateCommand.Applications, pipelineUpdateCommand.Stages, pipelineUpdateCommand.ThirdPartyServices)

	err := c.PipelineService.UpdatePipeline(pipelineUpdateCommand)
	if err != nil {

		switch err.(type) {
		case *models.CDPipelineDoesNotExistError:
			flash.Error(fmt.Sprintf("cd pipeline %v doesn't exist", pipelineName))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/%s/update", pipelineName), http.StatusFound)
			return
		case *models.NonValidRelatedBranchError:
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", pipelineUpdateCommand.Applications))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/%s/update", pipelineName), http.StatusBadRequest)
			return
		default:
			c.Abort("500")
			return
		}
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Redirect("/admin/edp/cd-pipeline/overview#cdPipelineSuccessModal", 302)
}

func (c *CDPipelineController) CreateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := c.GetStrings("app")
	pipelineName := c.GetString("pipelineName")
	serviceCheckboxes := c.GetStrings("service")
	stages := retrieveStagesFromRequest(c)

	cdPipelineCreateCommand := models.CDPipelineCommand{
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

		switch pipelineErr.(type) {
		case *models.CDPipelineExistsError:
			flash.Error(fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name))
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/cd-pipeline/create", http.StatusFound)
			return
		case *models.NonValidRelatedBranchError:
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications))
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/cd-pipeline/create", http.StatusBadRequest)
			return
		default:
			c.Abort("500")
			return
		}

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
		stageRequest := models.StageCreate{
			Name:        stageName,
			Description: this.GetString(stageName + "-stageDesc"),
			TriggerType: this.GetString(stageName + "-triggerType"),
			Order:       index,
		}

		for _, stepName := range this.GetStrings(stageName + "-stageStepName") {
			qualityGateRequest := models.QualityGate{
				QualityGateType: this.GetString(stageName + "-" + stepName + "-stageQualityGateType"),
				StepName:        stepName,
			}

			if qualityGateRequest.QualityGateType == "autotests" {
				autotestName := this.GetString(stageName + "-" + stepName + "-stageAutotests")
				qualityGateRequest.AutotestName = &autotestName
				stageName := this.GetString(stageName + "-" + stepName + "-stageBranch")
				qualityGateRequest.BranchName = &stageName
			}

			stageRequest.QualityGates = append(stageRequest.QualityGates, qualityGateRequest)
		}

		stages = append(stages, stageRequest)
	}

	log.Printf("Stages are fetched from request: %v", stages)
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
			Status: "inactive",
		}
		cdPipelines = append(cdPipelines, &pipeline)
	}
	return cdPipelines
}

func validateCDPipelineUpdateRequestData(cdPipeline models.CDPipelineCommand) *ErrMsg {
	isApplicationsValid := true
	isCDPipelineValid := true
	isStagesValid := true
	isQualityGatesValid := true
	errMsg := &ErrMsg{"An internal error has occurred on server while validating CD Pipeline's request body.", http.StatusInternalServerError}
	valid := validation.Validation{}
	isCDPipelineValid, err := valid.Valid(cdPipeline)

	if err != nil {
		return errMsg
	}

	if cdPipeline.Applications != nil {
		for _, app := range cdPipeline.Applications {
			isApplicationsValid, err = valid.Valid(app)
			if err != nil {
				return errMsg
			}
		}
	}

	if cdPipeline.Stages != nil {
		for _, stage := range cdPipeline.Stages {

			isValid, err := validateQualityGates(valid, stage.QualityGates)
			if err != nil {
				return errMsg
			}
			isQualityGatesValid = isValid

			isStagesValid, err = valid.Valid(stage)
			if err != nil {
				return errMsg
			}
		}
	}

	if isCDPipelineValid && isApplicationsValid && isStagesValid && isQualityGatesValid {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}

func validateQualityGates(valid validation.Validation, qualityGates []models.QualityGate) (bool, error) {
	isQualityGatesValid := true

	if qualityGates != nil {
		for _, qualityGate := range qualityGates {
			isValid, err := valid.Valid(qualityGate)
			if err != nil {
				return false, err
			}
			isQualityGatesValid = isValid

			if (qualityGate.QualityGateType == "autotests" && (qualityGate.AutotestName == nil || qualityGate.BranchName == nil)) ||
				(qualityGate.QualityGateType == "manual" && (qualityGate.AutotestName != nil || qualityGate.BranchName != nil)) {
				isQualityGatesValid = false
			}
		}
	} else {
		valid.Errors = append(valid.Errors, &validation.Error{Key: "qualityGates", Message: "can not be empty"})
		isQualityGatesValid = false
	}

	return isQualityGatesValid, nil
}
