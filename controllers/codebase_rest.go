/*
 * Copyright 2020 EPAM Systems.
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
	"edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	dberror "edp-admin-console/util/error/db-errors"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/microcosm-cc/bluemonday"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"path"
	"strings"
)

type CodebaseRestController struct {
	beego.Controller
	CodebaseService service.CodebaseService
}

func (c *CodebaseRestController) Prepare() {
	c.EnableXSRF = false
}

func (c *CodebaseRestController) GetCodebases() {
	criteria := c.getFilterCriteria()

	codebases, err := c.CodebaseService.GetCodebasesByCriteria(*criteria)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Data["json"] = codebases
	c.ServeJSON()
}

func (c *CodebaseRestController) getFilterCriteria() *query.CodebaseCriteria {
	return &query.CodebaseCriteria{
		Type:      c.getType(),
		Codebases: c.getCodebases(),
	}
}

func (c *CodebaseRestController) getType() *query.CodebaseType {
	codebaseType := c.GetString("type")
	if codebaseType == "" || validation.IsCodebaseTypeAcceptable(codebaseType) {
		cType := query.CodebaseTypes[codebaseType]
		return &cType
	}
	return nil
}

func (c *CodebaseRestController) getCodebases() []string {
	codebaseName := c.GetString("codebases")
	if codebaseName != "" {
		return strings.Split(codebaseName, ",")
	}
	return nil
}

func (c *CodebaseRestController) GetCodebase() {
	codebaseName := c.GetString(":codebaseName")
	if !bluemonday.UserSelectHandler(codebaseName) {
		http.Error(c.Ctx.ResponseWriter, "Incorrect CodebaseName", http.StatusInternalServerError)
		return
	}
	codebase, err := c.CodebaseService.GetCodebaseByName(codebaseName)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if codebase == nil {
		nonAppMsg := fmt.Sprintf("Please check codebase name. It seems there're not %s codebase.", codebaseName)
		http.Error(c.Ctx.ResponseWriter, html.EscapeString(nonAppMsg), http.StatusNotFound)
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

	errMsg := validation.ValidCodebaseRequestData(codebase)
	if errMsg != nil {
		log.Error("Failed to validate request data", zap.String("err", errMsg.Message))
		http.Error(c.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	ld := validation.CreateCodebaseLogRequestData(codebase)
	log.Info(ld.String())

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)
	if err != nil {
		c.checkError(err, codebase.Name, codebase.GitUrlPath)
		return
	}

	log.Info("Codebase resource is saved into cluster", zap.String("codebase", createdObject.Name))

	location := fmt.Sprintf("%s/%s", c.Ctx.Input.URL(), uuid.NewV4().String())
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Ctx.Output.Header("Location", location)
}

func (c *CodebaseRestController) checkError(err error, name string, url *string) {
	switch err.(type) {
	case *edperror.CodebaseAlreadyExistsError:
		errMsg := fmt.Sprintf("Codebase %v already exists.", name)
		http.Error(c.Ctx.ResponseWriter, html.EscapeString(errMsg), http.StatusBadRequest)
	case *edperror.CodebaseWithGitUrlPathAlreadyExistsError:
		errMsg := fmt.Sprintf("Codebase %v with %v project path already exists.", name, *url)
		http.Error(c.Ctx.ResponseWriter, html.EscapeString(errMsg), http.StatusBadRequest)
	default:
		log.Error("couldn't create codebase", zap.Error(err))
		errMsg := fmt.Sprintf("Failed to create codebase: %v", err.Error())
		http.Error(c.Ctx.ResponseWriter, errMsg, http.StatusInternalServerError)
	}
}

func (c *CodebaseRestController) Delete() {
	var cr command.DeleteCodebaseCommand
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&cr)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("delete codebase method is invoked", zap.String("codebase name", cr.Name))

	cdb, err := c.CodebaseService.GetCodebaseByName(cr.Name)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if cdb == nil {
		msg := fmt.Sprintf("Please check codebase name. It seems there's no %s codebase.", cr.Name)
		http.Error(c.Ctx.ResponseWriter, html.EscapeString(msg), http.StatusNotFound)
		return
	}

	if err := c.CodebaseService.Delete(cr.Name, string(cdb.Type)); err != nil {
		if dberror.CodebaseIsUsed(err) {
			cerr := err.(dberror.CodebaseIsUsedByCDPipeline)
			log.Error(cerr.Message, zap.Error(err))
			http.Error(c.Ctx.ResponseWriter, html.EscapeString(cerr.Message), http.StatusConflict)
			return
		}
		log.Error("delete process is failed", zap.Error(err))
		http.Error(c.Ctx.ResponseWriter, "delete process is failed", http.StatusInternalServerError)
		return
	}
	log.Info("delete codebase method is finished", zap.String("codebase name", cr.Name))

	location := fmt.Sprintf("%s/%s", c.Ctx.Input.URL(), uuid.NewV4().String())
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Ctx.Output.Header("Location", location)
}
