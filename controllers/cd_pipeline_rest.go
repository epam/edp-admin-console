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
	stage := models.Stage{
		Name:        "sit",
		CDPipeline:  "team-a",
		Description: "SIT environment for dedicated team",
		QualityGate: "manual",
		TriggerType: "is-changed",
		Applications: []models.ApplicationStage{
			{
				Name:     "petclinic-fe",
				InputIs:  "petclinic-fe-release-1.0",
				OutputIs: "team-a-sit-petclinic-fe-verified",
			},
			{
				Name:     "petclinic-be",
				InputIs:  "petclinic-be-master",
				OutputIs: "team-a-sit-petclinic-be-verified",
			},
		},
	}

	this.Data["json"] = stage
	this.ServeJSON()
}
