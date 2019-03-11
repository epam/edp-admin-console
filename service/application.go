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
	"fmt"
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
)

type ApplicationService struct {
	Clients               k8s.ClientSet
	ApplicationRepository repository.ApplicationEntityRepository
}

func (this ApplicationService) CreateApp(app models.App, edpName string) (*k8s.BusinessApplication, error) {
	appClient := this.Clients.ApplicationClient
	coreClient := this.Clients.CoreClient
	spec := convertData(app)
	namespace := edpName + "-edp-cicd"

	crd := &k8s.BusinessApplication{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "BusinessApplication",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: namespace,
		},
		Spec: spec,
		Status: k8s.BusinessApplicationStatus{
			Status: "initialized",
		},
	}

	err := createTempSecrets(namespace, app, coreClient)

	if err != nil {
		return nil, err
	}

	result := &k8s.BusinessApplication{}
	err = appClient.Post().Namespace(namespace).Resource("businessapplications").Body(crd).Do().Into(result)

	if err != nil {
		log.Printf("An error has occurred while creating object in k8s: %s", err)
		return &k8s.BusinessApplication{}, err
	}

	return result, nil
}

func (this ApplicationService) GetApplicationCR(appName string, edpName string) (*k8s.BusinessApplication, error) {
	appClient := this.Clients.ApplicationClient
	namespace := edpName + "-edp-cicd"

	result := &k8s.BusinessApplication{}
	err := appClient.Get().Namespace(namespace).Resource("businessapplications").Name(appName).Do().Into(result)

	if k8serrors.IsNotFound(err) {
		log.Printf("Current resourse %s doesn't exist.", appName)
		return nil, nil
	}

	if err != nil {
		log.Printf("An error has occurred while getting object from k8s: %s", err)
		return nil, err
	}

	return result, nil
}

func (this ApplicationService) GetAllApplications(edpName string) ([]models.BusinessEntity, error) {
	applications, err := this.ApplicationRepository.GetAllApplications(edpName)
	if err != nil {
		log.Printf("An error has occurred while getting application objects from database: %s", err)
		return nil, err
	}
	return applications, nil
}

func (this ApplicationService) GetApplication(appName string, edpName string) (*models.BusinessEntity, error) {
	application, err := this.ApplicationRepository.GetApplication(appName, edpName)
	if err != nil {
		log.Printf("An error has occurred while getting application object %s from database: %s", appName, err)
		return nil, err
	}

	return application, nil
}

func createSecret(namespace string, secret *v1.Secret, coreClient *coreV1Client.CoreV1Client) (*v1.Secret, error) {
	createdSecret, err := coreClient.Secrets(namespace).Create(secret)
	if err != nil {
		log.Printf("An error has occurred while saving secret: %s", err)
		return &v1.Secret{}, err
	}
	return createdSecret, nil
}

func createTempSecrets(namespace string, app models.App, coreClient *coreV1Client.CoreV1Client) error {
	if app.Repository != nil && (app.Repository.Login != "" && app.Repository.Password != "") {
		repoSecretName := fmt.Sprintf("repository-application-%s-temp", app.Name)
		tempRepoSecret := getSecret(repoSecretName, app.Repository.Login, app.Repository.Password)
		repositorySecret, err := createSecret(namespace, tempRepoSecret, coreClient)

		if err != nil {
			log.Printf("An error has occurred while creating repository secret: %s", err)
			return err
		}
		log.Printf("Repository secret was created: %s", repositorySecret)
	}

	if app.Vcs != nil {
		vcsSecretName := fmt.Sprintf("vcs-autouser-application-%s-temp", app.Name)
		tempVcsSecret := getSecret(vcsSecretName, app.Vcs.Login, app.Vcs.Password)
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

func convertData(app models.App) k8s.BusinessApplicationSpec {
	spec := k8s.BusinessApplicationSpec{
		Lang:      app.Lang,
		Framework: app.Framework,
		BuildTool: app.BuildTool,
		Strategy:  app.Strategy,
		Name:      app.Name,
	}

	if app.MultiModule {
		spec.Framework = app.Framework + "-multimodule"
	}

	if app.Repository != nil {
		spec.Repository = &k8s.Repository{
			Url: app.Repository.Url,
		}
	}

	if app.Route != nil {
		spec.Route = &k8s.Route{
			Site: app.Route.Site,
			Path: app.Route.Path,
		}
	}

	if app.Database != nil {
		spec.Database = &k8s.Database{
			Kind:     app.Database.Kind,
			Version:  app.Database.Version,
			Capacity: app.Database.Capacity,
			Storage:  app.Database.Storage,
		}
	}

	return spec
}
