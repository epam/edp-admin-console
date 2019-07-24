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
	"edp-admin-console/models/dto"
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type ICDPipelineRepository interface {
	GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error)
	GetCDPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error)
	GetStage(cdPipelineName, stageName string) (*models.StageView, error)
	GetCodebaseAndBranchName(codebaseId, branchId int) (*dto.CodebaseBranchDTO, error)
	GetQualityGates(stageId int64) ([]query.QualityGate, error)
}

const (
	SelectStageByCDPipelineAndStageNames = "select cs.name stage_name, " +
		"	cp.name pipeline_name, " +
		"	cs.description, " +
		"	cs.trigger_type, " +
		"	c.name app_name, " +
		"	cb.name branch_name, " +
		"	in_cds.oc_image_stream_name input_image_stream, " +
		"	out_cds.oc_image_stream_name output_image_stream, " +
		"	cs.\"order\", " +
		"   cs.id " +
		"from cd_stage cs " +
		"	left join stage_codebase_docker_stream scds on cs.id = scds.cd_stage_id " +
		"	left join codebase_docker_stream in_cds on scds.input_codebase_docker_stream_id = in_cds.id " +
		"	left join codebase_docker_stream out_cds on scds.output_codebase_docker_stream_id = out_cds.id " +
		"	left join codebase c on in_cds.codebase_id = c.id " +
		"	left join cd_pipeline_codebase_branch cpcb on cs.cd_pipeline_id = cpcb.cd_pipeline_id " +
		"	left join codebase_branch cb on cb.id = cpcb.codebase_branch_id " +
		"	left join cd_pipeline cp on cs.cd_pipeline_id = cp.id " +
		"where cp.name = ? " +
		"  and cs.name = ? " +
		"  and cb.codebase_id = in_cds.codebase_id;"
	SelectCodebaseAndBranchName = "select c.name codebase_name, cb.name codebase_branch_name " +
		"	from codebase c " +
		"left join codebase_branch cb on c.id = cb.codebase_id " +
		"where c.type = 'autotests' " +
		"  and c.id = ? " +
		"and cb.id = ? ;"
)

type CDPipelineRepository struct {
	ICDPipelineRepository
}

func (this CDPipelineRepository) GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error) {
	o := orm.NewOrm()
	cdPipeline := query.CDPipeline{Name: pipelineName}

	err := o.Read(&cdPipeline, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&cdPipeline, "CodebaseBranch", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	for _, c := range cdPipeline.CodebaseBranch {
		_, err = o.LoadRelated(c, "Codebase", false, 100, 0, "Name")
		if err != nil {
			return nil, err
		}
	}

	_, err = o.LoadRelated(&cdPipeline, "Stage", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	err = loadRelatedQualityGates(cdPipeline.Stage)
	if err != nil {
		return nil, err
	}

	for _, stage := range cdPipeline.Stage {
		err := loadRelatedAutotest(stage.QualityGates)
		if err != nil {
			return nil, err
		}

		err = loadRelatedBranch(stage.QualityGates)
		if err != nil {
			return nil, err
		}
	}

	_, err = o.LoadRelated(&cdPipeline, "ThirdPartyService", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	return &cdPipeline, nil
}

func loadRelatedQualityGates(stages []*query.Stage) error {
	for i, stage := range stages {
		o := orm.NewOrm()

		_, err := o.QueryTable(new(query.QualityGate)).
			Filter("cd_stage_id", stage.Id).
			All(&stage.QualityGates)
		if err != nil {
			return err
		}

		stages[i] = stage
	}

	return nil
}

func loadRelatedAutotest(gates []query.QualityGate) error {
	for i, gate := range gates {
		if gate.QualityGateType == "autotests" {
			o := orm.NewOrm()

			codebase := query.Codebase{}
			err := o.QueryTable(new(query.Codebase)).
				Filter("id", gate.CodebaseId).
				One(&codebase)
			if err != nil {
				return err
			}

			gates[i].Autotest = &codebase
		}
	}

	return nil
}

func loadRelatedBranch(gates []query.QualityGate) error {
	for i, gate := range gates {
		if gate.QualityGateType == "autotests" {
			o := orm.NewOrm()

			branch := query.CodebaseBranch{}
			err := o.QueryTable(new(query.CodebaseBranch)).
				Filter("id", gate.CodebaseBranchId).
				One(&branch)
			if err != nil {
				return err
			}

			gates[i].Branch = &branch
		}
	}

	return nil
}

func (this CDPipelineRepository) GetCDPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error) {
	o := orm.NewOrm()
	var pipelines []*query.CDPipeline

	qs := o.QueryTable(new(query.CDPipeline))

	if criteria.Status != "" {
		qs = qs.Filter("status", criteria.Status)
	}

	_, err := qs.OrderBy("name").All(&pipelines)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pipelines, nil
}

func (this CDPipelineRepository) GetStage(cdPipelineName, stageName string) (*models.StageView, error) {
	o := orm.NewOrm()
	var stage models.StageView
	var maps []orm.Params

	_, err := o.Raw(SelectStageByCDPipelineAndStageNames, cdPipelineName, stageName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for index, row := range maps {
		if index == 0 {
			stage.Name = row["stage_name"].(string)
			stage.CDPipeline = row["pipeline_name"].(string)
			stage.Description = row["description"].(string)
			stage.TriggerType = row["trigger_type"].(string)
			stage.Order = row["order"].(string)

			id, err := strconv.ParseInt(row["id"].(string), 10, 64)
			if err != nil {
				return nil, err
			}

			stage.Id = id
		}
		stage.Applications = append(stage.Applications, models.ApplicationStage{
			Name:       row["app_name"].(string),
			BranchName: row["branch_name"].(string),
			InputIs:    row["input_image_stream"].(string),
			OutputIs:   row["output_image_stream"].(string),
		})
	}

	return &stage, nil
}

func (CDPipelineRepository) GetCodebaseAndBranchName(codebaseId, branchId int) (*dto.CodebaseBranchDTO, error) {
	o := orm.NewOrm()

	result := dto.CodebaseBranchDTO{}
	var maps []orm.Params

	_, err := o.Raw(SelectCodebaseAndBranchName, codebaseId, branchId).Values(&maps)
	if err != nil {
		return nil, err
	}

	for _, row := range maps {
		result.AppName = row["codebase_name"].(string)
		result.BranchName = row["codebase_branch_name"].(string)
	}

	return &result, nil
}

func (CDPipelineRepository) GetQualityGates(stageId int64) ([]query.QualityGate, error) {
	o := orm.NewOrm()

	var gates []query.QualityGate
	_, err := o.QueryTable(new(query.QualityGate)).
		Filter("cd_stage_id", stageId).
		All(&gates)
	if err != nil {
		return nil, err
	}

	err = loadRelatedAutotest(gates)
	if err != nil {
		return nil, err
	}

	err = loadRelatedBranch(gates)
	if err != nil {
		return nil, err
	}

	return gates, nil
}
