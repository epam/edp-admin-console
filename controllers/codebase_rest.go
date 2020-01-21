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
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/satori/go.uuid"
	"net/http"
	"path"
)

type CodebaseRestController struct {
	beego.Controller
	CodebaseService service.CodebaseService
}

type ErrMsg struct {
	Message    string
	StatusCode int
}

func (c *CodebaseRestController) Prepare() {
	c.EnableXSRF = false
}

func (c *CodebaseRestController) GetCodebases() {
	criteria, err := getFilterCriteria(c)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	codebases, err := c.CodebaseService.GetCodebasesByCriteria(*criteria)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Data["json"] = codebases
	c.ServeJSON()
}

func isTypeAcceptable(getParam string) bool {
	if _, ok := query.CodebaseTypes[getParam]; ok {
		return true
	}
	return false
}

func getFilterCriteria(this *CodebaseRestController) (*query.CodebaseCriteria, error) {
	codebaseType := this.GetString("type")
	if codebaseType == "" || isTypeAcceptable(codebaseType) {
		return &query.CodebaseCriteria{
			Type: query.CodebaseTypes[codebaseType],
		}, nil
	}
	return nil, errors.New("type is not valid")
}

func (c *CodebaseRestController) GetCodebase() {
	codebaseName := c.GetString(":codebaseName")
	codebase, err := c.CodebaseService.GetCodebaseByName(codebaseName)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if codebase == nil {
		nonAppMsg := fmt.Sprintf("Please check codebase name. It seems there're not %s codebase.", codebaseName)
		http.Error(c.Ctx.ResponseWriter, nonAppMsg, http.StatusNotFound)
		return
	}

	c.Data["json"] = codebase
	c.ServeJSON()
}

func (c *CodebaseRestController) CreateCodebase() {
	var codebase command.CreateCodebase
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&codebase)
	usr, _ := c.Ctx.Input.Session("username").(string)
	codebase.Username = usr
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if codebase.Strategy != "import" {
		codebase.GitServer = "gerrit"
	} else {
		codebase.Name = path.Base(*codebase.GitUrlPath)
	}

	errMsg := ValidRequestData(codebase)
	if errMsg != nil {
		log.Info("Failed to validate request data", "err", errMsg.Message)
		http.Error(c.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	LogRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			errMsg := fmt.Sprintf("Codebase resource with %s name is already exists.", codebase.Name)
			http.Error(c.Ctx.ResponseWriter, errMsg, http.StatusBadRequest)
			return
		}
		errMsg := fmt.Sprintf("Failed to create codebase resource: %v", err.Error())
		http.Error(c.Ctx.ResponseWriter, errMsg, http.StatusInternalServerError)
		return
	}

	log.Info("Codebase resource is saved into cluster", "codebase", createdObject.Name)

	location := fmt.Sprintf("%s/%s", c.Ctx.Input.URL(), uuid.NewV4().String())
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Ctx.Output.Header("Location", location)
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
