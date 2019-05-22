package controllers

import (
	"edp-admin-console/models"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
)

type AutotestController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
	CodebaseService  service.CodebaseService
}

const paramWaitingForAutotest = "waitingforautotest"
const AutotestType = "autotest"

func (this *AutotestController) CreateAutotest() {
	flash := beego.NewFlash()
	codebase := extractAutotestRequestData(this)
	errMsg := validateAutotestRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate autotest request data: %s", errMsg.Message)
		flash := beego.NewFlash()
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/autotest/create", 302)
		return
	}
	log.Printf("Received autotest data from request: %+v %+v", codebase, codebase.Repository)

	createdObject, err := this.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Autotest %s is already exists.", codebase.Name)
			flash.Store(&this.Controller)
			this.Redirect("/admin/edp/autotest/create", 302)
			return
		}
		this.Abort("500")
		return
	}

	log.Printf("Autotest object is saved into k8s: %s", createdObject)
	flash.Success("Autotest object is created.")
	flash.Store(&this.Controller)
	this.Redirect(fmt.Sprintf("/admin/edp/autotest/overview?%s=%s#autotestSuccessModal", paramWaitingForAutotest, codebase.Name), 302)
}

func extractAutotestRequestData(this *AutotestController) models.Codebase {
	codebase := models.Codebase{
		Name:      this.GetString("nameOfApp"),
		Lang:      this.GetString("appLang"),
		Framework: this.GetString("framework"),
		BuildTool: this.GetString("buildTool"),
		Strategy:  "Clone",
		Type:      AutotestType,
	}

	testReportFramework := this.GetString("testReportFramework")
	if testReportFramework != "" {
		codebase.TestReportFramework = &testReportFramework
	}

	repoUrl := this.GetString("gitRepoUrl")
	if repoUrl != "" {
		codebase.Repository = &models.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := this.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			codebase.Repository.Login = this.GetString("repoLogin")
			codebase.Repository.Password = this.GetString("repoPassword")
		}
	}
	return codebase
}

func validateAutotestRequestData(autotest models.Codebase) *ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(autotest)

	if autotest.Repository != nil {
		_, err = valid.Valid(autotest.Repository)

		isAvailable := util.IsGitRepoAvailable(autotest.Repository.Url, autotest.Repository.Login, autotest.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}
	}

	if err != nil {
		return &ErrMsg{"An internal error has occurred on server while validating autotest's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}

func (this *AutotestController) GetCreateAutotestPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	if flash.Data["error"] != "" {
		this.Data["Error"] = flash.Data["error"]
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.TplName = "create_autotest.html"
}

func (this *AutotestController) GetAutotestsOverviewPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	if flash.Data["success"] != "" {
		this.Data["Success"] = true
	}

	var autotestType = "autotests"
	codebases, err := this.CodebaseService.GetAllCodebases(models.CodebaseCriteria{
		Type: &autotestType,
	})
	codebases = addCodebaseInProgressIfAny(codebases, this.GetString(paramWaitingForAutotest))
	if err != nil {
		this.Abort("500")
		return
	}

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["Autotests"] = codebases
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.TplName = "autotest.html"
}
