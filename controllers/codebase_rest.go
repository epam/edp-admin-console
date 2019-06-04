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
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CodebaseRestController struct {
	beego.Controller
	CodebaseService service.CodebaseService
}

type ErrMsg struct {
	Message    string
	StatusCode int
}

var codebaseTypes = map[string]string{
	"application": "",
	"autotests":   "",
	"library":     "",
}

func (this *CodebaseRestController) GetCodebases() {
	criteria, err := getFilterCriteria(this)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	codebases, err := this.CodebaseService.GetAllCodebases(*criteria)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = codebases
	this.ServeJSON()
}

func isTypeAcceptable(getParam string) bool {
	if _, ok := codebaseTypes[getParam]; ok {
		return true
	}
	return false
}

func getFilterCriteria(this *CodebaseRestController) (*models.CodebaseCriteria, error) {
	codebaseType := this.GetString("type")
	if codebaseType == "" || isTypeAcceptable(codebaseType) {
		return &models.CodebaseCriteria{
			Type: &codebaseType,
		}, nil
	}
	return nil, errors.New("type is not valid")
}

func (this *CodebaseRestController) GetCodebase() {
	codebaseName := this.GetString(":codebaseName")
	codebase, err := this.CodebaseService.GetCodebase(codebaseName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if codebase == nil {
		nonAppMsg := fmt.Sprintf("Please check codebase name. It seems there're not %s codebase.", codebaseName)
		http.Error(this.Ctx.ResponseWriter, nonAppMsg, http.StatusNotFound)
		return
	}

	this.Data["json"] = codebase
	this.ServeJSON()
}

func (this *CodebaseRestController) CreateCodebase() {
	var codebase models.Codebase
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&codebase)
	usr, _ := this.Ctx.Input.Session("username").(string)
	codebase.Username = usr
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := validRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	logRequestData(codebase)

	createdObject, err := this.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			errMsg := fmt.Sprintf("Codebase resource with %s name is already exists.", codebase.Name)
			http.Error(this.Ctx.ResponseWriter, errMsg, http.StatusBadRequest)
			return
		}
		errMsg := fmt.Sprintf("Failed to create codebase resource: %v", err.Error())
		http.Error(this.Ctx.ResponseWriter, errMsg, http.StatusInternalServerError)
		return
	}

	log.Printf("Codebase resource is saved into k8s: %+v", createdObject)

	location := fmt.Sprintf("%s/%s", this.Ctx.Input.URL(), uuid.NewV4().String())
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.Output.Header("Location", location)
}

func validRequestData(codebase models.Codebase) *ErrMsg {
	valid := validation.Validation{}
	var resErr error

	_, err := valid.Valid(codebase)
	resErr = err

	if codebase.Repository != nil {
		_, err := valid.Valid(codebase.Repository)

		isAvailable := util.IsGitRepoAvailable(codebase.Repository.Url, codebase.Repository.Login, codebase.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}

		resErr = err
	}

	if codebase.Route != nil {
		if len(codebase.Route.Path) > 0 {
			_, err := valid.Valid(codebase.Route)
			resErr = err
		} else {
			valid.Match(codebase.Route.Site, regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]$"), "Route.Site.Match")
		}
	}

	if codebase.Vcs != nil {
		_, err := valid.Valid(codebase.Vcs)
		resErr = err
	}

	if codebase.Database != nil {
		_, err := valid.Valid(codebase.Database)
		resErr = err
	}

	if !isTypeAcceptable(codebase.Type) {
		err := &validation.Error{Key: "repository", Message: "codebase type should be: application, autotests  or library"}
		valid.Errors = append(valid.Errors, err)
	}

	if codebase.Type == "autotests" && codebase.Strategy != "clone" {
		err := &validation.Error{Key: "repository", Message: "strategy for autotests must be 'clone'"}
		valid.Errors = append(valid.Errors, err)
	}

	if codebase.Type == "autotests" && codebase.Repository == nil {
		err := &validation.Error{Key: "repository", Message: "repository for autotests can't be null"}
		valid.Errors = append(valid.Errors, err)
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

func logRequestData(app models.Codebase) {
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
