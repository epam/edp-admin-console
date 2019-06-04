package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
	"strings"
)

type AutotestController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    service.BranchService
}

const AutotestType = "autotests"

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
	logAutotestRequestData(codebase)

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
	this.Redirect(fmt.Sprintf("/admin/edp/autotest/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func logAutotestRequestData(autotest models.Codebase) {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Request data to create codebase CR is valid. name=%s, strategy=%s, lang=%s, buildTool=%s, testReportFramework=%s",
		autotest.Name, autotest.Strategy, autotest.Lang, autotest.BuildTool, *autotest.TestReportFramework))

	if autotest.Repository != nil {
		result.WriteString(fmt.Sprintf(", repositoryUrl=%s, repositoryLogin=%s", autotest.Repository.Url, autotest.Repository.Login))
	}

	if autotest.Vcs != nil {
		result.WriteString(fmt.Sprintf(", vcsLogin=%s", autotest.Vcs.Login))
	}

	log.Println(result.String())
}

func extractAutotestRequestData(this *AutotestController) models.Codebase {
	codebase := models.Codebase{
		Name:      this.GetString("nameOfApp"),
		Lang:      this.GetString("appLang"),
		BuildTool: this.GetString("buildTool"),
		Strategy:  "clone",
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

	description := this.GetString("description")
	if description != "" {
		codebase.Description = &description
	}

	vcsLogin := this.GetString("vcsLogin")
	vcsPassword := this.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		codebase.Vcs = &models.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}
	codebase.Username = this.Ctx.Input.Session("username").(string)
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

	if autotest.Description == nil {
		err := &validation.Error{Key: "description", Message: "Description field can't be empty."}
		valid.Errors = append(valid.Errors, err)
	}

	if autotest.Vcs != nil {
		_, err = valid.Valid(autotest.Vcs)
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

	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue()
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.Data["IsVcsEnabled"] = isVcsEnabled
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
	codebases = addCodebaseInProgressIfAny(codebases, this.GetString(paramWaitingForCodebase))
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["Codebases"] = codebases
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.Data["Type"] = autotestType
	this.TplName = "codebase.html"
}

func (this *AutotestController) GetAutotestOverviewPage() {
	testName := this.GetString(":testName")
	codebases, err := this.CodebaseService.GetCodebase(testName)
	if err != nil {
		this.Abort("500")
		return
	}

	branchEntities, err := this.BranchService.GetAllReleaseBranchesByAppName(testName)
	branchEntities = addCodebaseBranchInProgressIfAny(branchEntities, this.GetString(paramWaitingForBranch))
	if err != nil {
		this.Abort("500")
		return
	}

	this.Data["ReleaseBranches"] = branchEntities
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["Autotests"] = codebases
	this.TplName = "codebase_overview.html"
}
