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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type ApplicationService struct {
	Clients k8s.ClientSet
}

func (this ApplicationService) CreateApp(app models.App, edpName string) (k8s.BusinessApplication, error) {
	appClient := this.Clients.ApplicationClient
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
		Spec:   spec,
		Status: k8s.BusinessApplicationStatus{},
	}

	result := &k8s.BusinessApplication{}
	err := appClient.Post().Namespace(namespace).Resource("businessapplications").Body(crd).Do().Into(result)

	if err != nil {
		log.Printf("An error has occurred during creating object in k8s: %s", err)
		return k8s.BusinessApplication{}, err
	}

	return *result, nil
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
