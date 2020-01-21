package controllers

import (
	"edp-admin-console/models/command"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego/validation"
	"net/http"
	"regexp"
	"strings"
)

func ValidRequestData(codebase command.CreateCodebase) *ErrMsg {
	valid := validation.Validation{}
	var resErr error

	_, err := valid.Valid(codebase)
	resErr = err

	if codebase.Strategy == "import" {
		valid.Match(codebase.GitUrlPath, regexp.MustCompile("^\\/.*$"), "Spec.GitUrlPath")
	}

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

func LogRequestData(app command.CreateCodebase) {
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

	log.Info(result.String())
}
