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
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"errors"
	"fmt"
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"time"
)

type CodebaseService struct {
	Clients             k8s.ClientSet
	ICodebaseRepository repository.ICodebaseRepository
	BranchService       CodebaseBranchService
}

const (
	CodebaseKind   = "Codebase"
	CodebasePlural = "codebases"
)

func (s CodebaseService) CreateCodebase(codebase command.CreateCodebase) (*k8s.Codebase, error) {
	log.Printf("Start creating Codebase resource: %v ...", codebase)

	codebaseCr, err := s.GetCodebaseCR(codebase.Name)
	if err != nil {
		log.Printf("An error has occurred while fetching Codebase CR from k8s: %s", codebase.Name)
		return nil, err
	}

	if codebaseCr != nil {
		log.Printf("Codebase %s is already exists in k8s.", codebaseCr.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	codebaseDb, err := s.GetCodebaseByName(codebase.Name)
	if err != nil {
		log.Printf("An error has occurred while fetching Codebase entity from DB: %s", codebase.Name)
		return nil, err
	}

	if codebaseDb != nil {
		log.Printf("Codebase %s is already exists in DB.", codebaseDb.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	edpClient := s.Clients.EDPRestClient
	coreClient := s.Clients.CoreClient

	crd := &k8s.Codebase{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       CodebaseKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      codebase.Name,
			Namespace: context.Namespace,
		},
		Spec: convertData(codebase),
		Status: k8s.CodebaseStatus{
			Available:       false,
			LastTimeUpdated: time.Now(),
			Status:          "initialized",
			Username:        codebase.Username,
			Action:          "codebase_registration",
			Result:          "success",
			Value:           "inactive",
		},
	}
	log.Printf("CR was generated : %v. Waiting to save ...", crd)

	err = createTempSecrets(context.Namespace, codebase, coreClient)

	if err != nil {
		return nil, err
	}

	result := &k8s.Codebase{}
	err = edpClient.Post().Namespace(context.Namespace).Resource(CodebasePlural).Body(crd).Do().Into(result)

	if err != nil {
		log.Printf("An error has occurred while creating codebase resource in k8s: %s", err)
		return &k8s.Codebase{}, err
	}

	_, err = s.BranchService.CreateCodebaseBranch(command.CreateCodebaseBranch{
		Name:     "master",
		Username: codebase.Username,
	}, codebase.Name)
	if err != nil {
		log.Printf("Error has been occurred during the master branch creation: %v", err)
		return &k8s.Codebase{}, err
	}
	return result, nil
}

func (s CodebaseService) GetCodebaseCR(codebaseName string) (*k8s.Codebase, error) {
	edpClient := s.Clients.EDPRestClient

	result := &k8s.Codebase{}
	err := edpClient.Get().Namespace(context.Namespace).Resource(CodebasePlural).Name(codebaseName).Do().Into(result)

	if k8serrors.IsNotFound(err) {
		log.Printf("Current codebase resourse %s doesn't exist in k8s.", codebaseName)
		return nil, nil
	}

	if err != nil {
		log.Printf("An error has occurred while getting codebase object from k8s: %s", err)
		return nil, err
	}

	return result, nil
}

func (s *CodebaseService) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	codebases, err := s.ICodebaseRepository.GetCodebasesByCriteria(criteria)
	if err != nil {
		log.Printf("An error has occurred while getting codebase objects: %s", err)
		return nil, err
	}
	log.Printf("Fetched codebases. Count: %v. Values: %v", len(codebases), codebases)

	return codebases, nil
}

func (s CodebaseService) GetCodebaseByName(name string) (*query.Codebase, error) {
	codebase, err := s.ICodebaseRepository.GetCodebaseByName(name)
	if err != nil {
		log.Printf("An error has occurred while getting codebase object %s: %s", name, err)
		return nil, err
	}
	log.Printf("Fetched codebase info: %+v", codebase)

	if codebase != nil {
		createBranchLinks(codebase.CodebaseBranch, codebase.Name, context.Tenant)
	}

	return codebase, nil
}

func (s CodebaseService) ExistCodebaseAndBranch(cbName, brName string) bool {
	return s.ICodebaseRepository.ExistCodebaseAndBranch(cbName, brName)
}

func createSecret(namespace string, secret *v1.Secret, coreClient *coreV1Client.CoreV1Client) (*v1.Secret, error) {
	createdSecret, err := coreClient.Secrets(namespace).Create(secret)
	if err != nil {
		log.Printf("An error has occurred while saving secret: %s", err)
		return &v1.Secret{}, err
	}
	return createdSecret, nil
}

func createTempSecrets(namespace string, codebase command.CreateCodebase, coreClient *coreV1Client.CoreV1Client) error {
	if codebase.Repository != nil && (codebase.Repository.Login != "" && codebase.Repository.Password != "") {
		repoSecretName := fmt.Sprintf("repository-codebase-%s-temp", codebase.Name)
		tempRepoSecret := getSecret(repoSecretName, codebase.Repository.Login, codebase.Repository.Password)
		repositorySecret, err := createSecret(namespace, tempRepoSecret, coreClient)

		if err != nil {
			log.Printf("An error has occurred while creating repository secret: %s", err)
			return err
		}
		log.Printf("Repository secret was created: %s", repositorySecret)
	}

	if codebase.Vcs != nil {
		vcsSecretName := fmt.Sprintf("vcs-autouser-codebase-%s-temp", codebase.Name)
		tempVcsSecret := getSecret(vcsSecretName, codebase.Vcs.Login, codebase.Vcs.Password)
		vcsSecret, err := createSecret(namespace, tempVcsSecret, coreClient)

		if err != nil {
			log.Printf("An error has occurred while creating vcs secret: %s", err)
			return err
		}
		log.Printf("Vcs secret was created: %s", vcsSecret)
	}

	return nil
}

func getSecret(name string, username string, password string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: map[string]string{
			"username": username,
			"password": password,
		},
	}
}

func convertData(codebase command.CreateCodebase) k8s.CodebaseSpec {
	spec := k8s.CodebaseSpec{
		Lang:      codebase.Lang,
		Framework: codebase.Framework,
		BuildTool: codebase.BuildTool,
		Strategy:  codebase.Strategy,
		Name:      codebase.Name,
		Type:      codebase.Type,
	}

	if codebase.Framework != nil {
		spec.Framework = codebase.Framework
	}

	if codebase.Repository != nil {
		spec.Repository = &k8s.Repository{
			Url: codebase.Repository.Url,
		}
	}

	if codebase.Route != nil {
		spec.Route = &k8s.Route{
			Site: codebase.Route.Site,
		}
		if len(codebase.Route.Path) > 0 {
			spec.Route.Path = codebase.Route.Path
		}
	}

	if codebase.Database != nil {
		spec.Database = &k8s.Database{
			Kind:     codebase.Database.Kind,
			Version:  codebase.Database.Version,
			Capacity: codebase.Database.Capacity,
			Storage:  codebase.Database.Storage,
		}
	}

	if codebase.TestReportFramework != nil {
		spec.TestReportFramework = codebase.TestReportFramework
	}

	if codebase.Description != nil {
		spec.Description = codebase.Description
	}

	return spec
}

func (s CodebaseService) checkAppAndBranch(apps []models.ApplicationWithBranch) bool {
	for _, app := range apps {
		exist := s.ICodebaseRepository.ExistActiveCodebaseAndBranch(app.ApplicationName, app.BranchName)

		if !exist {
			return false
		}
	}
	return true
}

func (s CodebaseService) GetApplicationsToPromote(cdPipelineId int) ([]string, error) {
	appsToPromote, err := s.ICodebaseRepository.SelectApplicationToPromote(cdPipelineId)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred while fetching Ids of applications which shoould be promoted: %v", err)
	}
	return s.selectApplicationNames(appsToPromote)
}

func (s CodebaseService) selectApplicationNames(applicationsToPromote []*query.ApplicationsToPromote) ([]string, error) {
	var result []string
	for _, app := range applicationsToPromote {
		codebase, err := s.ICodebaseRepository.GetCodebaseById(app.CodebaseId)
		if err != nil {
			return nil, fmt.Errorf("an error has occurred while fetching Codebase by Id %v: %v", app.CodebaseId, err)
		}
		result = append(result, codebase.Name)
	}

	log.Printf("Fetched Application to promote: %v", result)

	return result, nil
}
