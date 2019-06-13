package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
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
	BranchService    service.CodebaseBranchService
}

func (c *LibraryController) GetLibraryListPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["success"] != "" {
		c.Data["Success"] = true
	}

	codebases, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.Library,
	})
	codebases = addCodebaseInProgressIfAny(codebases, c.GetString(paramWaitingForCodebase))
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["Codebases"] = codebases
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(c.GetSession("realm_roles").([]string))
	c.Data["Type"] = query.Library
	c.TplName = "codebase.html"
}

func (c *LibraryController) GetCreatePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	isVcsEnabled, err := c.EDPTenantService.GetVcsIntegrationValue()
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(c.GetSession("realm_roles").([]string))
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.Library
	c.TplName = "create_library.html"
}

func (c *LibraryController) Create() {
	flash := beego.NewFlash()
	codebase := c.extractLibraryRequestData()
	errMsg := validateLibraryRequestData(codebase)
	if errMsg != nil {
		log.Printf("Failed to validate library request data: %s", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect("/admin/edp/library/create", 302)
		return
	}
	logLibraryRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Library %s is already exists.", codebase.Name)
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/library/create", 302)
			return
		}
		c.Abort("500")
		return
	}

	log.Printf("Library object is saved into k8s: %s", createdObject)
	flash.Success("Library object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/admin/edp/library/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func (c *LibraryController) extractLibraryRequestData() command.CreateCodebase {
	library := command.CreateCodebase{
		Name:      c.GetString("nameOfApp"),
		Lang:      c.GetString("appLang"),
		BuildTool: c.GetString("buildTool"),
		Strategy:  c.GetString("strategy"),
		Type:      "library",
	}

	repoUrl := c.GetString("gitRepoUrl")
	if repoUrl != "" {
		library.Repository = &command.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := c.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			library.Repository.Login = c.GetString("repoLogin")
			library.Repository.Password = c.GetString("repoPassword")
		}
	}

	vcsLogin := c.GetString("vcsLogin")
	vcsPassword := c.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		library.Vcs = &command.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}
	library.Username = c.Ctx.Input.Session("username").(string)
	return library
}

func validateLibraryRequestData(library command.CreateCodebase) *ErrMsg {
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

func logLibraryRequestData(library command.CreateCodebase) {
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
