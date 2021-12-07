/*
 * Copyright 2021 EPAM Systems.
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

package stage

import (
	"edp-admin-console/context"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/service/cd_pipeline"
	"edp-admin-console/service/logger"
	"fmt"
	"github.com/astaxie/beego"
	"go.uber.org/zap"
	"net/http"
)

var log = logger.GetLogger()

type stageController struct {
	beego.Controller
	PipelineService cd_pipeline.CDPipelineService
}

func InitStageController(ps cd_pipeline.CDPipelineService) stageController {
	return stageController{
		PipelineService: ps,
	}
}

func (c stageController) UpdateCDStage() {
	flash := beego.NewFlash()
	tType := c.GetStrings("triggerType")
	pn := c.GetString(":pipelineName")
	stages, err := c.getCDPipelineStages(pn)
	if err != nil {
		log.Error("an error has occurred while getting stage of cd pipeline", zap.Error(err))
		c.Abort("500")
		return
	}

	for i, name := range stages {
		sc := command.CDStageCommand{
			Name:        name,
			TriggerType: tType[i],
		}
		log.Debug("Request data is received to edit CD pipeline stage",
			zap.String("cd-pipeline", pn),
			zap.String("stage", name),
			zap.Any("triggerType", sc.TriggerType))
		err = c.PipelineService.UpdatePipelineStage(sc, pn)
		if err != nil {
			switch err.(type) {
			case *edperror.CDPipelineStageDoesNotExistError:
				flash.Error(fmt.Sprintf("cd pipeline stage %v doesn't exist", name))
				flash.Store(&c.Controller)
				c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/%s/overview", context.BasePath, pn), http.StatusFound)
				return
			default:
				c.Abort("500")
				return
			}
		}
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/%s/overview", context.BasePath, pn), 302)
}

func (c stageController) getCDPipelineStages(pipelineName string) ([]string, error) {
	stages, err := c.PipelineService.GetCDPipelineStages(pipelineName)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred while getting the stages of cd-pipeline")
	}
	return stages, nil
}
