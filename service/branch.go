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

type BranchService struct {
	Clients                  k8s.ClientSet
	IReleaseBranchRepository repository.IReleaseBranchRepository
}

func (this *BranchService) CreateReleaseBranch(branchInfo models.ReleaseBranchCreateCommand, appName string) (*k8s.CodebaseBranch, error) {
	log.Println("Start creating CR for branch release...")
	edpRestClient := this.Clients.EDPRestClient

	releaseBranchCR, err := getReleaseBranchCR(edpRestClient, branchInfo.Name, appName, context.Namespace)
	if err != nil {
		log.Printf("An error has occurred while getting release branch CR {%s} from k8s: %s", branchInfo.Name, err)
		return nil, err
	}

	if releaseBranchCR != nil {
		log.Printf("release branch CR {%s} already exists in k8s", branchInfo.Name)
		return nil, errors.New(fmt.Sprintf("release branch CR {%s} already exists in k8s", branchInfo.Name))
	}

	spec := convertBranchInfoData(branchInfo, appName)
	branch := &k8s.CodebaseBranch{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "CodebaseBranch",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", appName, branchInfo.Name),
			Namespace: context.Namespace,
		},
		Spec: spec,
		Status: k8s.CodebaseBranchStatus{
			Status:          "initialized",
			LastTimeUpdated: time.Now(),
			Username:        branchInfo.Username,
			Action:          "codebase_branch_registration",
			Result:          "success",
		},
	}

	result := &k8s.CodebaseBranch{}
	err = edpRestClient.Post().Namespace(context.Namespace).Resource("codebasebranches").Body(branch).Do().Into(result)
	if err != nil {
		log.Printf("An error has occurred while creating release branch custom resource in k8s: %s", err)
		return &k8s.CodebaseBranch{}, err
	}
	return result, nil
}

func (this *BranchService) GetReleaseBranch(appName string, branchName string) (*models.ReleaseBranchView, error) {
	releaseBranch, err := this.IReleaseBranchRepository.GetReleaseBranch(appName, branchName)
	if err != nil {
		log.Printf("An error has occurred while getting branch entity %s-%s from database: %s", appName, branchName, err)
		return nil, err
	}

	if releaseBranch != nil {
		log.Printf("Fetched branch entity: {%+v}", releaseBranch)
	}

	return releaseBranch, nil
}

func (this *BranchService) GetAllReleaseBranchesByAppName(appName string) ([]models.ReleaseBranchView, error) {
	releaseBranches, err := this.IReleaseBranchRepository.GetAllReleaseBranchesByAppName(appName)
	if err != nil {
		log.Printf("An error has occurred while getting branch entities for {%s} application: %s", appName, err)
		return nil, err
	}

	if len(releaseBranches) != 0 {
		edpTenantName := context.Tenant
		createLinks(releaseBranches, appName, edpTenantName)
		log.Printf("Fetched branch entities: {%s}. Count: {%s}", releaseBranches, string(len(releaseBranches)))
	}

	return releaseBranches, nil
}

func (this *BranchService) GetAllReleaseBranches(branchFilterCriteria models.BranchCriteria) ([]models.ReleaseBranchView, error) {
	releaseBranches, err := this.IReleaseBranchRepository.GetAllReleaseBranches(branchFilterCriteria)
	if err != nil {
		log.Printf("An error has occurred while getting branch entities: %s", err)
		return nil, err
	}
	return releaseBranches, nil
}

func convertBranchInfoData(branchInfo models.ReleaseBranchCreateCommand, appName string) k8s.CodebaseBranchSpec {
	return k8s.CodebaseBranchSpec{
		Name:         branchInfo.Name,
		Commit:       branchInfo.Commit,
		CodebaseName: appName,
	}
}

func createLinks(branchEntities []models.ReleaseBranchView, appName string, edpTenantName string) {
	wildcard := beego.AppConfig.String("dnsWildcard")
	for index, branch := range branchEntities {
		branch.VCSLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/gitweb?p=%s.git;a=shortlog;h=refs/heads/%s", "gerrit", edpTenantName, wildcard, appName, branch.Name)
		branch.CICDLink = fmt.Sprintf("https://%s-%s-edp-cicd.%s/job/%s/view/%s", "jenkins", edpTenantName, wildcard, appName, strings.ToUpper(branch.Name))
		branchEntities[index] = branch
	}
}

func getReleaseBranchCR(edpRestClient *rest.RESTClient, branchName string, appName string, namespace string) (*k8s.CodebaseBranch, error) {
	result := &k8s.CodebaseBranch{}
	err := edpRestClient.Get().Namespace(namespace).Resource("codebasebranches").Name(fmt.Sprintf("%s-%s", appName, branchName)).Do().Into(result)

	if k8serrors.IsNotFound(err) {
		log.Printf("Current resourse %s doesn't exist.", branchName)
		return nil, nil
	}

	if err != nil {
		log.Printf("An error has occurred while getting release branch custom resource from k8s: %s", err)
		return nil, err
	}

	return result, nil
}
