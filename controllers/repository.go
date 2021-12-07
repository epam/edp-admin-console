package controllers

import (
	validation2 "edp-admin-console/controllers/validation"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type RepositoryRestController struct {
	beego.Controller
	AppService service.CodebaseService
}

type RepoData struct {
	Url      string `json:"url,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (this *RepositoryRestController) IsGitRepoAvailable() {
	var repo RepoData
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&repo); err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validRepoRequestData(repo); err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Message, err.StatusCode)
		return
	}

	resp, err := getResponseJSON(util.IsGitRepoAvailable(repo.Url, repo.Login, repo.Password))
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = resp
	this.ServeJSON()
}

func validRepoRequestData(repo RepoData) *validation2.ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(repo)
	if err != nil {
		return &validation2.ErrMsg{Message: "An error has occurred while validating application's form fields.", StatusCode: http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &validation2.ErrMsg{Message: string(validation2.CreateErrorResponseBody(valid)), StatusCode: http.StatusBadRequest}
}

func getResponseJSON(available bool, err error) (*string, error) {
	buf, err := json.Marshal(struct {
		Available bool   `json:"available"`
		Msg       string `json:"msg"`
	}{
		Available: available,
		Msg:       err.Error(),
	})
	if err != nil {
		return nil, err
	}
	return util.GetStringP(string(buf)), nil
}
