package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
)

type LibraryController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

const LibraryType = "library"

func (this *LibraryController) GetLibraryListPage() {
	this.Data["Codebases"] = []models.CodebaseView{
		{
			Name:      "stub-name-01",
			Language:  "java",
			BuildTool: "maven",
			Status:    "active",
		},
		{
			Name:      "stub-name-02",
			Language:  "java",
			BuildTool: "maven",
			Status:    "active",
		},
	}
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.Data["Type"] = LibraryType
	this.TplName = "codebase.html"
}

func (this *LibraryController) GetCreatePage() {
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
