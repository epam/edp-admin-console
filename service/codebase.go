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
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"time"
)

type CodebaseService struct {
	Clients             k8s.ClientSet
	ICodebaseRepository repository.ICodebaseEntityRepository
	BranchService       BranchService
}

const (
	CodebaseKind   = "Codebase"
	CodebasePlural = "codebases"
)

func (this CodebaseService) CreateCodebase(codebase models.Codebase) (*k8s.Codebase, error) {
	log.Printf("Start creating Codebase resource: %v ...", codebase)

	codebaseCr, err := this.GetCodebaseCR(codebase.Name)
	if err != nil {
		log.Printf("An error has occurred while fetching Codebase CR from k8s: %s", codebase.Name)
		return nil, err
	}

	codebaseDb, err := this.GetCodebase(codebase.Name)
	if err != nil {
		log.Printf("An error has occurred while fetching Codebase entity from DB: %s", codebase.Name)
		return nil, err
	}

	if codebaseCr != nil {
		log.Printf("Codebase %s is already exists in k8s.", codebaseCr.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	if codebaseDb != nil {
		log.Printf("Codebase %s is already exists in DB.", codebaseDb.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	edpClient := this.Clients.EDPRestClient
	coreClient := this.Clients.CoreClient

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

	_, err = this.BranchService.CreateReleaseBranch(models.ReleaseBranchCreateCommand{
		Name:     "master",
		Username: codebase.Username,
	}, codebase.Name)
	if err != nil {
		log.Printf("Error has been occurred during the master branch creation: %v", err)
		return &k8s.Codebase{}, err
	}
	return result, nil
}

func (this CodebaseService) GetCodebaseCR(codebaseName string) (*k8s.Codebase, error) {
	edpClient := this.Clients.EDPRestClient

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

func (this *CodebaseService) GetAllCodebases(criteria models.CodebaseCriteria) ([]models.CodebaseView, error) {
	codebases, err := this.ICodebaseRepository.GetAllCodebases(criteria)
	if err != nil {
		log.Printf("An error has occurred while getting codebase objects from database: %s", err)
		return nil, err
	}
	log.Printf("Fetched codebases. Count: %v. Rows: %v", len(codebases), codebases)

	return codebases, nil
}

func (this CodebaseService) GetCodebase(codebaseName string) (*models.CodebaseDetailInfo, error) {
	codebase, err := this.ICodebaseRepository.GetCodebase(codebaseName)
	if err != nil {
		log.Printf("An error has occurred while getting codebase object %s from database: %s", codebaseName, err)
		return nil, err
	}
	log.Printf("Fetched codebase info: %+v", codebase)

	return codebase, nil
}

func (this CodebaseService) GetAllCodebasesWithReleaseBranches(criteria models.CodebaseCriteria) ([]models.CodebaseWithReleaseBranch, error) {
	codebases, err := this.ICodebaseRepository.GetAllCodebasesWithReleaseBranches(criteria)
	if err != nil {
		log.Printf("An error has occurred while getting all codebase objects: %s", err)
		return nil, err
	}
	log.Printf("Fetched codebases info: {%+v}", codebases)

	return codebases, nil
}

func createSecret(namespace string, secret *v1.Secret, coreClient *coreV1Client.CoreV1Client) (*v1.Secret, error) {
	createdSecret, err := coreClient.Secrets(namespace).Create(secret)
	if err != nil {
		log.Printf("An error has occurred while saving secret: %s", err)
		return &v1.Secret{}, err
	}
	return createdSecret, nil
}

func createTempSecrets(namespace string, codebase models.Codebase, coreClient *coreV1Client.CoreV1Client) error {
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

func convertData(codebase models.Codebase) k8s.CodebaseSpec {
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

func (this CodebaseService) GetCodebaseByCodebaseAndBranchNames(codebaseName, branchName string) (*models.CodebaseView, error) {
	codebase, err := this.ICodebaseRepository.GetCodebaseByCodebaseAndBranchNames(codebaseName, branchName)
	if err != nil {
		log.Printf("An error has occurred while getting codebase object %s from database: %s", codebaseName, err)
		return nil, err
	}
	log.Printf("Fetched codebase info: %+v", codebase)

	return codebase, nil
}
