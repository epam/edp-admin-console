package controllers

import (
	"edp-admin-console/models"
	"edp-admin-console/service"
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

	log.Printf("Received data from client to create crd: %s", app)

	errMsg := validRequestData(err, app)
	if errMsg != (ErrMsg{}) {
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, errMsg.StatusCode)
		return
	}

	id := uuid.NewV4().String()

	createdObject, err := this.AppService.CreateApp(app)

	if err != nil {
		log.Println("Failed to create custom resource: " + err.Error())
		http.Error(this.Ctx.ResponseWriter, "Failed to create custom resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Custom object is saved into k8s: %s", createdObject)

	location := fmt.Sprintf("%s/%s", this.Ctx.Input.URL(), id)
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.Output.Header("Location", location)
}

func validRequestData(err error, addApp models.App) ErrMsg {
	valid := validation.Validation{}
	b, err := valid.Valid(addApp)
	if err != nil {
		return ErrMsg{err.Error(), http.StatusInternalServerError}
	}
	if !b {
		return ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
	}
	return ErrMsg{}
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
