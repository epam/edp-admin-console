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
	"time"
)

type CDPipelineService struct {
	Clients               k8s.ClientSet
	ICDPipelineRepository repository.ICDPipelineRepository
}

func (this *CDPipelineService) CreatePipeline(pipelineName string, releaseBranchCommands []models.ReleaseBranchCreatePipelineCommand) (*k8s.CDPipeline, error) {
	log.Println("Start creating CR pipeline...")
	edpRestClient := this.Clients.EDPRestClient
	namespace := beego.AppConfig.String("cicdNamespace") + "-edp-cicd"

	pipelineCR, err := getCDPipelineCR(edpRestClient, pipelineName, namespace)
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
			Namespace: namespace,
		},
		Spec: convertPipelineData(pipelineName, releaseBranchCommands),
		Status: k8s.CDPipelineStatus{
			LastTimeUpdated: time.Now().Format("2006-01-02T15:04:05Z"),
			Status:          "initialized",
		},
	}

	cdPipeline := &k8s.CDPipeline{}
	err = edpRestClient.Post().Namespace(namespace).Resource("cdpipelines").Body(crd).Do().Into(cdPipeline)

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
	log.Printf("Fetched CD Pipeline from DB: %s", cdPipeline)
	return cdPipeline, nil
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
