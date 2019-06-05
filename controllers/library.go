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

type LibraryController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
	CodebaseService  service.CodebaseService
	BranchService    service.BranchService
}

const LibraryType = "library"

func (this *LibraryController) GetLibraryListPage() {
	flash := beego.ReadFromRequest(&this.Controller)
	if flash.Data["success"] != "" {
		this.Data["Success"] = true
	}

	var libraryType = "library"
	codebases, err := this.CodebaseService.GetAllCodebases(models.CodebaseCriteria{
		Type: &libraryType,
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
	this.Data["Type"] = LibraryType
	this.TplName = "codebase.html"
}

func (this *LibraryController) GetCreatePage() {
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
	this.TplName = "create_library.html"
}

func (this *LibraryController) Create() {
	flash := beego.NewFlash()
	codebase := extractLibraryRequestData(this)
	errMsg := validateLibraryRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate library request data: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&this.Controller)
		this.Redirect("/admin/edp/library/create", 302)
		return
	}
	logLibraryRequestData(codebase)

	createdObject, err := this.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Library %s is already exists.", codebase.Name)
			flash.Store(&this.Controller)
			this.Redirect("/admin/edp/library/create", 302)
			return
		}
		this.Abort("500")
		return
	}

	log.Printf("Library object is saved into k8s: %s", createdObject)
	flash.Success("Library object is created.")
	flash.Store(&this.Controller)
	this.Redirect(fmt.Sprintf("/admin/edp/library/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func extractLibraryRequestData(this *LibraryController) models.Codebase {
	library := models.Codebase{
		Name:      this.GetString("nameOfApp"),
		Lang:      this.GetString("appLang"),
		BuildTool: this.GetString("buildTool"),
		Strategy:  this.GetString("strategy"),
		Type:      LibraryType,
	}

	repoUrl := this.GetString("gitRepoUrl")
	if repoUrl != "" {
		library.Repository = &models.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := this.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			library.Repository.Login = this.GetString("repoLogin")
			library.Repository.Password = this.GetString("repoPassword")
		}
	}

	vcsLogin := this.GetString("vcsLogin")
	vcsPassword := this.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		library.Vcs = &models.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}
	library.Username = this.Ctx.Input.Session("username").(string)
	return library
}

func validateLibraryRequestData(library models.Codebase) *ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(library)

	if library.Strategy == "clone" && library.Repository != nil {
		_, err = valid.Valid(library.Repository)

		isAvailable := util.IsGitRepoAvailable(library.Repository.Url, library.Repository.Login, library.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}
	}

	if library.Vcs != nil {
		_, err = valid.Valid(library.Vcs)
	}

	if err != nil {
		return &ErrMsg{"An internal error has occurred on server while validating autotest's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}

func logLibraryRequestData(library models.Codebase) {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Request data to create codebase CR is valid. name=%s, strategy=%s, lang=%s, buildTool=%s",
		library.Name, library.Strategy, library.Lang, library.BuildTool))

	if library.Repository != nil {
		result.WriteString(fmt.Sprintf(", repositoryUrl=%s, repositoryLogin=%s", library.Repository.Url, library.Repository.Login))
	}

	if library.Vcs != nil {
		result.WriteString(fmt.Sprintf(", vcsLogin=%s", library.Vcs.Login))
	}

	log.Println(result.String())
}
