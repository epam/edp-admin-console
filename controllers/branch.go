package controllers

import (
	"edp-admin-console/models"
	"edp-admin-console/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
	"regexp"
)

type BranchController struct {
	beego.Controller
	BranchService service.BranchService
}

func (this *BranchController) CreateReleaseBranch() {
	branchInfo := extractReleaseBranchRequestData(this)
	errMsg := validReleaseBranchRequestData(branchInfo)
	appName := this.GetString(":appName")
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		this.Redirect(fmt.Sprintf("/admin/edp/application/%s/overview", appName), 302)
		return
	}
	log.Printf("Request data to create CR for release branch is valid: {Branch Name: %s, Commit Hash: %s}", branchInfo.Name, branchInfo.Commit)

	branch, err := this.BranchService.GetReleaseBranch(appName, branchInfo.Name)
	if err != nil {
		this.Abort("500")
		return
	}

	if branch != nil {
		this.Redirect(fmt.Sprintf("/admin/edp/application/%s/overview#branchExistsModal", appName), 302)
		return
	}

	releaseBranch, err := this.BranchService.CreateReleaseBranch(branchInfo, appName)
	if err != nil {
		this.Abort("500")
		return
	}

	log.Printf("BranchRelease resource is saved into k8s: %s", releaseBranch)
	this.Redirect(fmt.Sprintf("/admin/edp/application/%s/overview#branchSuccessModal", appName), 302)
}

func extractReleaseBranchRequestData(this *BranchController) models.ReleaseBranchRequestData {
	return models.ReleaseBranchRequestData{
		Name:   this.GetString("name"),
		Commit: this.GetString("commit"),
	}
}

func validReleaseBranchRequestData(requestData models.ReleaseBranchRequestData) *ErrMsg {
	valid := validation.Validation{}
	_, err := valid.Valid(requestData)

	if len(requestData.Commit) != 0 {
		valid.Match(requestData.Commit, regexp.MustCompile("\\b([a-f0-9]{40})\\b"), "Commit.Match")
	}

	if err != nil {
		return &ErrMsg{"An internal error has occurred on server while validating branch's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}
