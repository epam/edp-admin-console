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

package controllers

import (
	"edp-admin-console/models"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ApplicationRestController struct {
	beego.Controller
	AppService service.ApplicationService
}

type ErrMsg struct {
	Message    string
	StatusCode int
}

func (this *ApplicationRestController) GetApplications() {
	edpTenantName := this.GetString(":name")
	applications, err := this.AppService.GetAllApplications(edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = applications
	this.ServeJSON()
}

func (this *ApplicationRestController) GetApplication() {
	edpTenantName := this.GetString(":name")
	appName := this.GetString(":appName")
	application, err := this.AppService.GetApplication(appName, edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if application == nil {
		nonAppMsg := fmt.Sprintf("Please check application name. It seems there're not %s application.", appName)
		http.Error(this.Ctx.ResponseWriter, nonAppMsg, http.StatusNotFound)
		return
	}

	this.Data["json"] = application
	this.ServeJSON()
}

func (this *ApplicationRestController) CreateApplication() {
	edpTenantName := this.GetString(":name")
	var app models.App
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&app)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := validRequestData(app)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	logRequestData(app)

	applicationCr, err := this.AppService.GetApplicationCR(app.Name, edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, "Failed to get custom resource from cluster: "+err.Error(), http.StatusInternalServerError)
		return
	}

	application, err := this.AppService.GetApplication(app.Name, edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, "Failed to get custom resource from database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if applicationCr != nil || application != nil {
		http.Error(this.Ctx.ResponseWriter, "Application name is already exists.", http.StatusBadRequest)
		return
	}

	createdObject, err := this.AppService.CreateApp(app, edpTenantName)

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, "Failed to create custom resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Custom object is saved into k8s: %s", createdObject)

	location := fmt.Sprintf("%s/%s", this.Ctx.Input.URL(), uuid.NewV4().String())
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.Output.Header("Location", location)
}

func validRequestData(addApp models.App) *ErrMsg {
	valid := validation.Validation{}
	var resErr error

	_, err := valid.Valid(addApp)
	resErr = err

	if addApp.Repository != nil {
		_, err := valid.Valid(addApp.Repository)

		isAvailable := util.IsGitRepoAvailable(addApp.Repository.Url, addApp.Repository.Login, addApp.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}

		resErr = err
	}

	if addApp.Route != nil {
		if len(addApp.Route.Path) > 0 {
			_, err := valid.Valid(addApp.Route)
			resErr = err
		} else {
			valid.Match(addApp.Route.Site, regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]$"), "Route.Site.Match")
		}
	}

	if addApp.Vcs != nil {
		_, err := valid.Valid(addApp.Vcs)
		resErr = err
	}

	if addApp.Database != nil {
		_, err := valid.Valid(addApp.Database)
		resErr = err
	}

	if resErr != nil {
		return &ErrMsg{"An internal error has occurred on server while validating application's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}

func createErrorResponseBody(valid validation.Validation) []byte {
	errJson, _ := json.Marshal(extractErrors(valid))
	errResponse := struct {
		Message string
		Content string
	}{
		"Body of request are not valid.",
		string(errJson),
	}
	response, _ := json.Marshal(errResponse)
	return response
}

func extractErrors(valid validation.Validation) []string {
	var errMap []string
	for _, err := range valid.Errors {
		errMap = append(errMap, fmt.Sprintf("Validation failed on %s: %s", err.Key, err.Message))
	}
	return errMap
}

func logRequestData(app models.App) {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Request data to create CR is valid. name=%s, strategy=%s, lang=%s, buildTool=%s, multiModule=%s, framework=%s",
		app.Name, app.Strategy, app.Lang, app.BuildTool, app.MultiModule, app.Framework))

	if app.Repository != nil {
		result.WriteString(fmt.Sprintf(", repositoryUrl=%s, repositoryLogin=%s", app.Repository.Url, app.Repository.Login))
	}

	if app.Vcs != nil {
		result.WriteString(fmt.Sprintf(", vcsLogin=%s", app.Vcs.Login))
	}

	if app.Route != nil {
		result.WriteString(fmt.Sprintf(", routeSite=%s, routePath=%s", app.Route.Site, app.Route.Path))
	}

	if app.Database != nil {
		result.WriteString(fmt.Sprintf(", dbKind=%s, dbМersion=%s, dbCapacity=%s, dbStorage=%s", app.Database.Kind, app.Database.Version, app.Database.Capacity, app.Database.Storage))
	}

	log.Println(result.String())
}