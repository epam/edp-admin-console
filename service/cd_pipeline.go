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
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"fmt"
	"github.com/astaxie/beego"
	appsV1Client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"log"
	"sort"
	"strings"
	"time"
)

type CDPipelineService struct {
	Clients               k8s.ClientSet
	ICDPipelineRepository repository.ICDPipelineRepository
	CodebaseService       CodebaseService
	BranchService         CodebaseBranchService
}

type ErrMsg struct {
	Message    string
	StatusCode int
}

const OpenshiftProjectLink = "%s/console/project/"

func (s *CDPipelineService) CreatePipeline(cdPipeline models.CDPipelineCommand) (*k8s.CDPipeline, error) {
	log.Printf("Start creating CD Pipeline: %v", cdPipeline)

	exist := s.CodebaseService.checkAppAndBranch(cdPipeline.Applications)
	if !exist {
		return nil, models.NewNonValidRelatedBranchError()
	}

	cdPipelineReadModel, err := s.GetCDPipelineByName(cdPipeline.Name)
	if err != nil {
		return nil, err
	}

	if cdPipelineReadModel != nil {
		log.Printf("CD Pipeline %s is already exists in DB.", cdPipeline.Name)
		return nil, models.NewCDPipelineExistsError()
	}

	edpRestClient := s.Clients.EDPRestClient
	pipelineCR, err := s.getCDPipelineCR(cdPipeline.Name)
	if err != nil {
		return nil, err
	}

	if pipelineCR != nil {
		log.Printf("CD Pipeline %s is already exists in k8s.", cdPipeline.Name)
		return nil, models.NewCDPipelineExistsError()
	}

	crd := &k8s.CDPipeline{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "CDPipeline",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipeline.Name,
			Namespace: context.Namespace,
		},
		Spec: convertPipelineData(cdPipeline),
		Status: k8s.CDPipelineStatus{
			LastTimeUpdated: time.Now(),
			Status:          "initialized",
		},
	}

	cdPipelineCr := &k8s.CDPipeline{}
	err = edpRestClient.Post().Namespace(context.Namespace).Resource("cdpipelines").Body(crd).Do().Into(cdPipelineCr)

	if err != nil {
		log.Printf("An error has occurred while creating CD Pipeline object in k8s: %s", err)
		return nil, err
	}
	log.Printf("Pipeline CR %v is saved into k8s", cdPipeline)

	_, err = s.CreateStages(edpRestClient, cdPipeline)
	if err != nil {
		log.Printf("An error has occurred while creating Stages in k8s: %s", err)
		return nil, err
	}
	log.Printf("Stages for CD Pipeline %s were created in k8s: %v", cdPipeline.Name, cdPipeline.Stages)

	return cdPipelineCr, nil
}

func (s *CDPipelineService) GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error) {
	log.Println("Start execution of GetCDPipelineByName method...")
	cdPipeline, err := s.ICDPipelineRepository.GetCDPipelineByName(pipelineName)
	if err != nil {
		log.Printf("An error has occurred while getting CD Pipeline from database: %s", err)
		return nil, err
	}
	if cdPipeline != nil {
		createJenkinsLink(cdPipeline)
		if len(cdPipeline.Stage) != 0 {
			sortStagesByOrder(cdPipeline.Stage)
			createOpenshiftProjectLinks(cdPipeline.Stage, cdPipeline.Name)
			log.Printf("Fetched Stages. Count: {%v}. Rows: {%v}", len(cdPipeline.Stage), cdPipeline.Stage)
		}
		for i, branch := range cdPipeline.CodebaseBranch {
			branch.AppName = branch.Codebase.Name
			cdPipeline.CodebaseBranch[i] = branch
		}

		matrix, err := fillCodebaseStageMatrix(s.Clients.AppsV1Client, cdPipeline)
		if err == nil {
			cdPipeline.CodebaseStageMatrix = matrix
		}

		applicationsToPromote, err := s.CodebaseService.GetApplicationsToPromote(cdPipeline.Id)
		if err != nil {
			log.Printf("An error has occurred while getting Applications To Promote for CD Pipeline %v: %v", cdPipeline.Id, err)
			return nil, err
		}

		cdPipeline.ApplicationsToPromote = applicationsToPromote

		log.Printf("Fetched CD Pipeline from DB: %v", cdPipeline)
	}

	return cdPipeline, nil
}

func (s *CDPipelineService) CreateStages(edpRestClient *rest.RESTClient, cdPipeline models.CDPipelineCommand) ([]k8s.Stage, error) {
	log.Printf("Start creating CR stages: %+v\n", cdPipeline.Stages)

	err := checkStagesInK8s(edpRestClient, cdPipeline.Name, cdPipeline.Stages)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stagesCr, err := saveStagesIntoK8s(edpRestClient, cdPipeline.Name, cdPipeline.Stages)
	if err != nil {
		return nil, err
	}

	return stagesCr, nil
}

func (s *CDPipelineService) GetAllPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error) {
	log.Println("Start fetching all CD Pipelines...")
	cdPipelines, err := s.ICDPipelineRepository.GetCDPipelines(criteria)
	if err != nil {
		log.Printf("An error has occurred while getting CD Pipelines from database: %s", err)
		return nil, err
	}

	if len(cdPipelines) != 0 {
		createJenkinsLinks(cdPipelines)
	}
	log.Printf("Fetched CD Pipelines. Count: {%v}. Rows: {%v}", len(cdPipelines), cdPipelines)

	return cdPipelines, nil
}

func (s *CDPipelineService) UpdatePipeline(pipeline models.CDPipelineCommand) error {
	log.Printf("Start updating CD Pipeline: %v", pipeline.Name)

	if pipeline.Applications != nil {
		exist := s.CodebaseService.checkAppAndBranch(pipeline.Applications)
		if !exist {
			return models.NewNonValidRelatedBranchError()
		}
	}

	cdPipelineReadModel, err := s.GetCDPipelineByName(pipeline.Name)
	if err != nil {
		return err
	}

	if cdPipelineReadModel == nil {
		log.Printf("CD Pipeline %s doesn't exist in DB.", pipeline.Name)
		return models.NewCDPipelineDoesNotExistError()
	}

	pipelineCR, err := s.getCDPipelineCR(pipeline.Name)
	if err != nil {
		return err
	}

	if pipelineCR == nil {
		log.Printf("CD Pipeline %s doesn't exist in k8s.", pipeline.Name)
		return models.NewCDPipelineDoesNotExistError()
	}

	if pipeline.Applications != nil {
		log.Printf("Start updating Autotest for CD Pipeline: %v. New Applications: %v", pipelineCR.Spec.Name, pipeline.Applications)

		var codebaseBranches []string
		for _, v := range pipeline.Applications {
			codebaseBranches = append(codebaseBranches, fmt.Sprintf("%s-%s", v.ApplicationName, v.BranchName))
		}

		pipelineCR.Spec.CodebaseBranch = codebaseBranches
	}

	pipelineCR.Spec.ApplicationsToPromote = pipeline.ApplicationToApprove
	pipelineCR.Status.LastTimeUpdated = time.Now()

	edpRestClient := s.Clients.EDPRestClient

	err = edpRestClient.Put().
		Namespace(context.Namespace).
		Resource("cdpipelines").
		Name(pipelineCR.Spec.Name).
		Body(pipelineCR).
		Do().
		Into(pipelineCR)

	if err != nil {
		log.Printf("An error has occurred while updating CD Pipeline object in k8s: %s", err)
		return err
	}

	log.Printf("CD Pipeline %v has been updated with new data; %v", pipeline.Name, pipeline)

	return nil
}

func sortStagesByOrder(stages []*query.Stage) {
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Order < stages[j].Order
	})
}

func (s *CDPipelineService) GetStage(cdPipelineName, stageName string) (*models.StageView, error) {
	log.Printf("Start fetching Stage by CD Pipeline %s and Stage %s names...", cdPipelineName, stageName)
	stage, err := s.ICDPipelineRepository.GetStage(cdPipelineName, stageName)
	if err != nil {
		log.Printf("An error has occurred while getting Stage from database: %s", err)
		return nil, err
	}

	if stage == nil {
		log.Printf("Couldn't find Stage by %v name and %v CD Pipeline name", stageName, cdPipelineName)
		return nil, nil
	}

	gates, err := s.ICDPipelineRepository.GetQualityGates(stage.Id)
	if err != nil {
		log.Printf("An error has occurred while fetching Quality Gates from database: %v", err)
		return nil, err
	}
	stage.QualityGates = gates

	log.Printf("Fetched Stage: {%v}, Quality Gates: {%v}", stage, stage.QualityGates)

	return stage, nil
}

func createOpenshiftProjectLinks(stages []*query.Stage, cdPipelineName string) {
	for index, stage := range stages {
		stage.OpenshiftProjectName = fmt.Sprintf("%s-%s-%s", context.Tenant, cdPipelineName, stage.Name)
		stage.OpenshiftProjectLink = fmt.Sprintf(OpenshiftProjectLink+stage.OpenshiftProjectName, beego.AppConfig.String("openshiftClusterURL"))
		stages[index] = stage
	}
}

func fillCodebaseStageMatrix(ocClient *appsV1Client.AppsV1Client, cdPipeline *query.CDPipeline) (map[query.CDCodebaseStageMatrixKey]query.CDCodebaseStageMatrixValue, error) {
	var matrix = make(map[query.CDCodebaseStageMatrixKey]query.CDCodebaseStageMatrixValue, len(cdPipeline.CodebaseBranch)*len(cdPipeline.Stage))
	for _, stage := range cdPipeline.Stage {
		dcs, err := ocClient.DeploymentConfigs(stage.OpenshiftProjectName).List(metav1.ListOptions{})
		if err != nil {
			log.Printf("An error has occurred while getting project from OpenShift: %s", err)
			return nil, err
		}
		for _, codebase := range cdPipeline.CodebaseBranch {
			var key = query.CDCodebaseStageMatrixKey{
				CodebaseBranch: codebase,
				Stage:          stage,
			}
			var value = query.CDCodebaseStageMatrixValue{
				DockerVersion: "no deploy",
			}
			for _, dc := range dcs.Items {
				for _, container := range dc.Spec.Template.Spec.Containers {
					if container.Name == codebase.AppName {
						var containerImage = container.Image
						var delimeter = strings.LastIndex(containerImage, ":")
						if delimeter > 0 {
							value.DockerVersion = string(containerImage[(delimeter + 1):len(containerImage)])
						}
					}
				}
			}
			matrix[key] = value
		}

	}
	return matrix, nil
}

func convertPipelineData(cdPipeline models.CDPipelineCommand) k8s.CDPipelineSpec {
	var codebaseBranches []string
	for _, v := range cdPipeline.Applications {
		codebaseBranches = append(codebaseBranches, fmt.Sprintf("%s-%s", v.ApplicationName, v.BranchName))
	}
	return k8s.CDPipelineSpec{
		Name:                  cdPipeline.Name,
		CodebaseBranch:        codebaseBranches,
		ThirdPartyServices:    cdPipeline.ThirdPartyServices,
		ApplicationsToPromote: cdPipeline.ApplicationToApprove,
	}
}

func (s *CDPipelineService) getCDPipelineCR(pipelineName string) (*k8s.CDPipeline, error) {
	edpRestClient := s.Clients.EDPRestClient
	cdPipeline := &k8s.CDPipeline{}

	err := edpRestClient.Get().Namespace(context.Namespace).Resource("cdpipelines").Name(pipelineName).Do().Into(cdPipeline)

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

func createJenkinsLinks(cdPipelines []*query.CDPipeline) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	for index, pipeline := range cdPipelines {
		pipeline.JenkinsLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s", "jenkins", context.Tenant, wildcard, fmt.Sprintf("%s-%s", pipeline.Name, "cd-pipeline"))
		cdPipelines[index] = pipeline
		log.Printf("Created Jenkins link %v", pipeline.JenkinsLink)
	}
}

func createJenkinsLink(cdPipeline *query.CDPipeline) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	cdPipeline.JenkinsLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s", "jenkins", context.Tenant, wildcard, fmt.Sprintf("%s-%s", cdPipeline.Name, "cd-pipeline"))
	log.Printf("Created CD Pipeline Jenkins link %v", cdPipeline.JenkinsLink)
	createLinksForBranchEntities(cdPipeline.CodebaseBranch)
}

func createLinksForBranchEntities(branchEntities []*query.CodebaseBranch) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	for index, branch := range branchEntities {
		branch.VCSLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/gitweb?p=%s.git;a=shortlog;h=refs/heads/%s", "gerrit", context.Tenant, wildcard, branch.Codebase.Name, branch.Name)
		branch.CICDLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s/view/%s", "jenkins", context.Tenant, wildcard, branch.Codebase.Name, strings.ToUpper(branch.Name))
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
			Name:         stage.Name,
			Description:  stage.Description,
			TriggerType:  stage.TriggerType,
			Order:        stage.Order,
			CdPipeline:   cdPipelineName,
			QualityGates: stage.QualityGates,
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
