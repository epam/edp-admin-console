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

package repository

import (
	"edp-admin-console/models"
	"edp-admin-console/repository/sql_builder"
	"github.com/astaxie/beego/orm"
)

type ICDPipelineRepository interface {
	GetCDPipelineByName(pipelineName string) (*models.CDPipelineDTO, error)
	GetCDPipelines(filterCriteria models.CDPipelineCriteria) ([]models.CDPipelineView, error)
	GetCDPipelineStages(pipelineName string) ([]models.CDPipelineStageView, error)
}

const (
	SelectCDPipelineByName = "select cdp.name as pipeline_name, cb.name as branch_name, c.name as app_name, cdp.status " +
		"from cd_pipeline cdp " +
		"		left join cd_pipeline_codebase_branch cpcb on cdp.id = cpcb.cd_pipeline_id " +
		"		left join codebase_branch cb on cpcb.codebase_branch_id = cb.id " +
		"		left join codebase c on cb.codebase_id = c.id " +
		"where cdp.name = ?;"
	SelectStagesByName = "select cs.name, cs.description, cs.trigger_type, cs.quality_gate, cs.jenkins_step_name " +
		"from cd_stage cs " +
		"		left join cd_pipeline cp on cs.cd_pipeline_id = cp.id " +
		"where cp.name = ?;"
)

type CDPipelineRepository struct {
	ICDPipelineRepository
}

func (this CDPipelineRepository) GetCDPipelineByName(pipelineName string) (*models.CDPipelineDTO, error) {
	o := orm.NewOrm()
	var cdPipeline models.CDPipelineDTO
	var maps []orm.Params

	_, err := o.Raw(SelectCDPipelineByName, pipelineName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for _, row := range maps {
		cdPipeline.Name = row["pipeline_name"].(string)
		cdPipeline.Status = row["status"].(string)
		cdPipeline.CodebaseBranches = append(cdPipeline.CodebaseBranches, models.CodebaseBranchDTO{
			AppName:    row["app_name"].(string),
			BranchName: row["branch_name"].(string),
		})
	}

	return &cdPipeline, nil
}

func (this CDPipelineRepository) GetCDPipelines(filterCriteria models.CDPipelineCriteria) ([]models.CDPipelineView, error) {
	o := orm.NewOrm()
	var pipelines []models.CDPipelineView

	selectAllCDPipelinesQuery := sql_builder.GetAllCDPipelinesQuery(filterCriteria)
	_, err := o.Raw(selectAllCDPipelinesQuery).QueryRows(&pipelines)

	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (this CDPipelineRepository) GetCDPipelineStages(pipelineName string) ([]models.CDPipelineStageView, error) {
	o := orm.NewOrm()
	var stages []models.CDPipelineStageView

	_, err := o.Raw(SelectStagesByName, pipelineName).QueryRows(&stages)

	if err != nil {
		return nil, err
	}

	return stages, nil
}
