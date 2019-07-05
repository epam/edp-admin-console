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
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
)

type CDPipelineRestController struct {
	beego.Controller
	CDPipelineService service.CDPipelineService
}

func (this *CDPipelineRestController) GetCDPipelineByName() {
	pipelineName := this.GetString(":name")
	cdPipeline, err := this.CDPipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if cdPipeline == nil {
		nonAppMsg := fmt.Sprintf("Please check CD Pipeline name. It seems there's not %s pipeline.", pipelineName)
		http.Error(this.Ctx.ResponseWriter, nonAppMsg, http.StatusNotFound)
		return
	}

	this.Data["json"] = cdPipeline
	this.ServeJSON()
}

func (this *CDPipelineRestController) GetStage() {
	pipelineName := this.GetString(":pipelineName")
	stageName := this.GetString(":stageName")

	stage, err := this.CDPipelineService.GetStage(pipelineName, stageName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if stage == nil {
		http.Error(this.Ctx.ResponseWriter, "Please check request data.", http.StatusNotFound)
		return
	}

	this.Data["json"] = stage
	this.ServeJSON()
}

func (this *CDPipelineRestController) CreateCDPipeline() {
	var cdPipelineCreateCommand models.CDPipelineCreateCommand
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&cdPipelineCreateCommand)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	errMsg := validateCDPipelineRequestData(cdPipelineCreateCommand)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	log.Printf("Request data is receieved to create CD pipeline: %s. Applications: %v. Stages: %v. Services: %v",
		cdPipelineCreateCommand.Name, cdPipelineCreateCommand.Applications, cdPipelineCreateCommand.Stages, cdPipelineCreateCommand.ThirdPartyServices)

	_, pipelineErr := this.CDPipelineService.CreatePipeline(cdPipelineCreateCommand)
	if pipelineErr != nil {

		switch pipelineErr.(type) {
		case *models.CDPipelineExistsError:
			http.Error(this.Ctx.ResponseWriter, fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name), http.StatusFound)
			return
		case *models.NonValidRelatedBranchError:
			http.Error(this.Ctx.ResponseWriter, fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications), http.StatusBadRequest)
			return
		default:
			http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	this.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

func validateCDPipelineRequestData(cdPipeline models.CDPipelineCreateCommand) *ErrMsg {
	var isCDPipelineValid, isApplicationsValid, isStagesValid bool
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

			if (stage.QualityGateType == "autotests" && stage.Autotests == nil) ||
				(stage.QualityGateType == "manual" && stage.Autotests != nil) {
				isStagesValid = false
			}

			isStagesValid, err = valid.Valid(stage)
			if err != nil {
				return errMsg
			}
		}
	}

	if isCDPipelineValid && isApplicationsValid && isStagesValid {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}
