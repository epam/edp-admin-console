package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/k8s"
	"edp-admin-console/util/consts"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/astaxie/beego"
	"github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type JiraRestController struct {
	beego.Controller
	Clients k8s.ClientSet
}

func (c *JiraRestController) GetJiraMetadataFields() {
	sn := c.GetString(":serverName")
	logv := log.With(zap.String("server name", sn))
	logv.Info("start GetJiraMetadataFields method")
	jc, err := c.initJiraClient(sn)
	if err != nil {
		log.Error(err.Error())
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	fields, err := c.getJiraMetadataFields(*jc)
	if err != nil {
		log.Error(err.Error())
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	logv.Info("end GetJiraMetadataFields method")
	c.Data["json"] = fields
	c.ServeJSON()
}

func (c *JiraRestController) initJiraClient(jiraServerName string) (*gojira.Client, error) {
	js, err := c.getJiraServerCr(jiraServerName)
	if err != nil {
		return nil, errors.Wrapf(err, "an error has occurred while getting %v Jira server CR", jiraServerName)
	}

	s, err := c.getSecret(js.Spec.CredentialName, js.Namespace)
	if err != nil {
		return nil, err
	}

	user := string(s.Data["username"])
	pwd := string(s.Data["password"])
	jc, err := initClient(user, pwd, js.Spec.ApiUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create Jira client. user - %v, apiUrl - %v", user, js.Spec.ApiUrl)
	}
	return jc, nil
}

func (c *JiraRestController) getJiraServerCr(name string) (*v1alpha1.JiraServer, error) {
	js := &v1alpha1.JiraServer{}
	err := c.Clients.EDPRestClient.
		Get().
		Namespace(context.Namespace).
		Resource(consts.JiraServerPlural).
		Name(name).
		Do().
		Into(js)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func (c *JiraRestController) getSecret(name, namespace string) (*v1.Secret, error) {
	s, err := c.Clients.CoreClient.Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func initClient(user, password, apiUrl string) (*gojira.Client, error) {
	tp := gojira.BasicAuthTransport{
		Username: user,
		Password: password,
	}
	client, err := gojira.NewClient(tp.Client(), apiUrl)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (c *JiraRestController) getJiraMetadataFields(client gojira.Client) (map[string]string, error) {
	log.Info("start fetching Jira metadata fields")
	respFields, _, err := client.Field.GetList()
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string, len(respFields))
	for _, v := range respFields {
		fields[v.ID] = v.Name
	}
	log.Info("fields have been fetched", zap.Int("count", len(fields)))
	return fields, nil
}
