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
)

type AppController struct {
	beego.Controller
	AppService service.ApplicationService
}

type ErrMsg struct {
	Message    string
	StatusCode int
}

func (this *AppController) CreateApplication() {
	var app models.App
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&app)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := validRequestData(app)
	if errMsg != nil {
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, errMsg.StatusCode)
		return
	}

	id := uuid.NewV4().String()

	edpTenantName := this.GetString(":name")
	createdObject, err := this.AppService.CreateApp(app, edpTenantName)

	if err != nil {
		log.Printf("Failed to create custom resource in %s namespace: %s", edpTenantName, err.Error())
		http.Error(this.Ctx.ResponseWriter, "Failed to create custom resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Custom object is saved into k8s: %s", createdObject)

	location := fmt.Sprintf("%s/%s", this.Ctx.Input.URL(), id)
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
		_, err := valid.Valid(addApp.Route)
		resErr = err
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
		return &ErrMsg{"An error has occurred while validating application's form fields.", http.StatusInternalServerError}
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
