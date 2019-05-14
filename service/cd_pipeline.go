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

package service

import (
	"edp-admin-console/context"
	"edp-admin-console/k8s"
	"edp-admin-console/models"
	"edp-admin-console/repository"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"log"
	"strings"
	"time"
)

type CDPipelineService struct {
	Clients               k8s.ClientSet
	ICDPipelineRepository repository.ICDPipelineRepository
}

const OpenshiftProjectLink = "https://master.delivery.aws.main.edp.projects.epam.com/console/project/"

func (this *CDPipelineService) CreatePipeline(pipelineName string, releaseBranchCommands []models.ReleaseBranchCreatePipelineCommand) (*k8s.CDPipeline, error) {
	log.Println("Start creating CR pipeline...")
	edpRestClient := this.Clients.EDPRestClient

	pipelineCR, err := getCDPipelineCR(edpRestClient, pipelineName, context.Namespace)
	if err != nil {
		return nil, err
	}

	if pipelineCR != nil {
		log.Printf("pipeline CR {%s} already exists in k8s", pipelineName)
		return nil, errors.New(fmt.Sprintf("pipeline CR {%s} already exists in k8s", pipelineName))
	}

	crd := &k8s.CDPipeline{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "CDPipeline",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pipelineName,
			Namespace: context.Namespace,
		},
		Spec: convertPipelineData(pipelineName, releaseBranchCommands),
		Status: k8s.CDPipelineStatus{
			LastTimeUpdated: time.Now(),
			Status:          "initialized",
		},
	}

	cdPipeline := &k8s.CDPipeline{}
	err = edpRestClient.Post().Namespace(context.Namespace).Resource("cdpipelines").Body(crd).Do().Into(cdPipeline)

	if err != nil {
		log.Printf("An error has occurred while creating CD Pipeline object in k8s: %s", err)
		return nil, err
	}
	log.Println("Pipeline CR is saved into k8s")
	return cdPipeline, nil
}

func (this *CDPipelineService) GetCDPipelineByName(pipelineName string) (*models.CDPipelineDTO, error) {
	log.Println("Start execution of GetCDPipelineByName method...")
	cdPipeline, err := this.ICDPipelineRepository.GetCDPipelineByName(pipelineName)
	if err != nil {
		log.Printf("An error has occurred while getting CD Pipeline from database: %s", err)
		return nil, err
	}
	if cdPipeline != nil {
		createJenkinsLink(cdPipeline)
		log.Printf("Fetched CD Pipeline from DB: %s", cdPipeline)
	}

	return cdPipeline, nil
}

func (this *CDPipelineService) CreateStages(cdPipelineName string, stages []models.StageCreate) ([]k8s.Stage, error) {
	log.Printf("Start creating CR stages: %+v\n", stages)
	edpRestClient := this.Clients.EDPRestClient

	err := checkStagesInK8s(edpRestClient, cdPipelineName, stages)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stagesCr, err := saveStagesIntoK8s(edpRestClient, cdPipelineName, stages)
	if err != nil {
		return nil, err
	}

	return stagesCr, nil
}

func (this *CDPipelineService) GetAllPipelines(filterCriteria models.CDPipelineCriteria) ([]models.CDPipelineView, error) {
	log.Println("Start fetching all CD Pipelines...")
	cdPipelines, err := this.ICDPipelineRepository.GetCDPipelines(filterCriteria)
	if err != nil {
		log.Printf("An error has occurred while getting CD Pipelines from database: %s", err)
		return nil, err
	}

	if len(cdPipelines) != 0 {
		createJenkinsLinks(cdPipelines)
		log.Printf("Fetched CD Pipelines. Count: {%v}. Rows: {%v}", len(cdPipelines), cdPipelines)
	}

	return cdPipelines, nil
}

func (this *CDPipelineService) GetCDPipelineStages(cdPipelineName string) ([]models.CDPipelineStageView, error) {
	log.Printf("Start fetching all stages for %v CD Pipeline...", cdPipelineName)
	stages, err := this.ICDPipelineRepository.GetCDPipelineStages(cdPipelineName)
	if err != nil {
		log.Printf("An error has occurred while getting Stages from database: %s", err)
		return nil, err
	}

	if len(stages) != 0 {
		createOpenshiftProjectLinks(stages, cdPipelineName)
		log.Printf("Fetched Stages. Count: {%v}. Rows: {%v}", len(stages), stages)
	}

	return stages, nil
}

func createOpenshiftProjectLinks(stages []models.CDPipelineStageView, cdPipelineName string) {
	for index, stage := range stages {
		stage.OpenshiftProjectLink = fmt.Sprintf(OpenshiftProjectLink+"%s-%s-%s", context.Tenant, cdPipelineName, stage.Name)
		stages[index] = stage
	}
}

func convertPipelineData(pipelineName string, releaseBranchCommands []models.ReleaseBranchCreatePipelineCommand) k8s.CDPipelineSpec {
	var codebaseBranches []string
	for _, v := range releaseBranchCommands {
		codebaseBranches = append(codebaseBranches, fmt.Sprintf("%s-%s", v.AppName, v.BranchName))
	}
	return k8s.CDPipelineSpec{
		Name:           pipelineName,
		CodebaseBranch: codebaseBranches,
	}
}

func getCDPipelineCR(edpRestClient *rest.RESTClient, pipelineName string, namespace string) (*k8s.CDPipeline, error) {
	cdPipeline := &k8s.CDPipeline{}
	err := edpRestClient.Get().Namespace(namespace).Resource("cdpipelines").Name(pipelineName).Do().Into(cdPipeline)

	if k8serrors.IsNotFound(err) {
		log.Printf("Pipeline resource %s doesn't exist.", pipelineName)
		return nil, nil
	}

	if err != nil {
		log.Printf("An error has occurred while getting Pipeline CR from k8s: %s", err)
		return nil, err
	}

	return cdPipeline, nil
}

func createJenkinsLinks(cdPipelines []models.CDPipelineView) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	for index, pipeline := range cdPipelines {
		pipeline.JenkinsLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s", "jenkins", context.Tenant, wildcard, fmt.Sprintf("%s-%s", pipeline.Name, "cd-pipeline"))
		cdPipelines[index] = pipeline
		log.Printf("Created Jenkins link %v", pipeline.JenkinsLink)
	}
}

func createJenkinsLink(cdPipeline *models.CDPipelineDTO) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	cdPipeline.JenkinsLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s", "jenkins", context.Tenant, wildcard, fmt.Sprintf("%s-%s", cdPipeline.Name, "cd-pipeline"))
	log.Printf("Created CD Pipeline Jenkins link %v", cdPipeline.JenkinsLink)
	createLinksForBranchEntities(cdPipeline.CodebaseBranches)
}

func createLinksForBranchEntities(branchEntities []models.CodebaseBranchDTO) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	for index, branch := range branchEntities {
		branch.BranchLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/gitweb?p=%s.git;a=shortlog;h=refs/heads/%s", "gerrit", context.Tenant, wildcard, branch.AppName, branch.BranchName)
		branch.JenkinsLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s/view/%s", "jenkins", context.Tenant, wildcard, branch.AppName, strings.ToUpper(branch.BranchName))
		branchEntities[index] = branch
	}
}

func createCrd(cdPipelineName string, stage models.StageCreate) k8s.Stage {
	return k8s.Stage{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "Stage",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", cdPipelineName, stage.Name),
			Namespace: context.Namespace,
		},
		Spec: k8s.StageSpec{
			Name:        stage.Name,
			Description: stage.Description,
			QualityGate: stage.QualityGateType,
			JenkinsStep: stage.StepName,
			TriggerType: stage.TriggerType,
			Order:       stage.Order,
			CdPipeline:  cdPipelineName,
		},
		Status: k8s.StageStatus{
			LastTimeUpdated: time.Now(),
			Status:          "initialized",
		},
	}
}

func saveStagesIntoK8s(edpRestClient *rest.RESTClient, cdPipelineName string, stages []models.StageCreate) ([]k8s.Stage, error) {
	var stagesCr []k8s.Stage
	for _, stage := range stages {
		crd := createCrd(cdPipelineName, stage)
		stageCr := k8s.Stage{}
		err := edpRestClient.Post().Namespace(context.Namespace).Resource("stages").Body(&crd).Do().Into(&stageCr)
		if err != nil {
			log.Printf("An error has occurred while creating Stage object in k8s: %s", err)
			return nil, err
		}
		log.Printf("Stage is saved into k8s: %+v\n", stage.Name)
		stagesCr = append(stagesCr, stageCr)
	}
	return stagesCr, nil
}

func checkStagesInK8s(edpRestClient *rest.RESTClient, cdPipelineName string, stages []models.StageCreate) error {
	for _, stage := range stages {
		stagesCr := &k8s.Stage{}
		stageName := fmt.Sprintf("%s-%s", cdPipelineName, stage.Name)
		err := edpRestClient.Get().Namespace(context.Namespace).Resource("stages").Name(stageName).Do().Into(stagesCr)

		if k8serrors.IsNotFound(err) {
			log.Printf("Stage %s doesn't exist.", stage.Name)
			continue
		}

		if err != nil {
			log.Printf("An error has occurred while getting Stage from k8s: %s", err)
			return err
		}

		if stagesCr != nil {
			log.Printf("CR Stage %s is already exists in k8s: %s", stageName, err)
			return fmt.Errorf("stage %s is already exists", stage.Name)
		}
	}
	return nil
}
