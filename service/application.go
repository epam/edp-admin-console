package service

import (
	"edp-admin-console/k8s"
	"edp-admin-console/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"log"
)

type ApplicationService struct {
	CrdClient *rest.RESTClient
}

func (this ApplicationService) CreateApp(app models.App) (k8s.BusinessApplication, error) {
	spec := convertData(app)

	crd := &k8s.BusinessApplication{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "edp.epam.com/v1alpha1",
			Kind:       "BusinessApplication",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: "anton-test-edp-k8s",
		},
		Spec:   spec,
		Status: k8s.BusinessApplicationStatus{},
	}

	result := &k8s.BusinessApplication{}
	err := this.CrdClient.Post().Namespace("anton-test-edp-k8s").Resource("businessapplications").Body(crd).Do().Into(result)

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

	if app.Git != nil {
		spec.Git = &k8s.Git{
			Url:      app.Git.Url,
			Login:    app.Git.Login,
			Password: app.Git.Password,
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
