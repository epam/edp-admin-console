package controllers

import (
	"edp-admin-console/models"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
)

type AutotestController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *AutotestController) CreateAutotest() {

}

func (this *AutotestController) GetCreateAutotestPage() {

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

	version, err := this.EDPTenantService.GetEDPVersion()
	if err != nil {
		this.Abort("500")
		return
	}

	autotests := []models.AutotestView{
		{
			Name:      "stub-autotest-1",
			Status:    "in_progress",
			BuildTool: "maven",
			Language:  "java",
		},
		{
			Name:      "stub-autotest-2",
			Status:    "created",
			BuildTool: "maven",
			Language:  "java",
		},
		{
			Name:      "stub-autotest-3",
			Status:    "failed",
			BuildTool: "maven",
			Language:  "java",
		},
	}

	this.Data["Autotests"] = autotests
	this.Data["EDPVersion"] = version
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["HasRights"] = isAdmin(this.GetSession("realm_roles").([]string))
	this.TplName = "autotest.html"
}
