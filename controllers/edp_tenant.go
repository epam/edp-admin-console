package controllers

import (
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"net/http"
)

type EDPTenantController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
}

func (this *EDPTenantController) GetEDPTenants() {
	resourceAccess := this.Ctx.Input.Session("resource_access").(map[string][]string)
	edpTenants, err := this.EDPTenantService.GetEDPTenants(resourceAccess)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["InputURL"] = this.Ctx.Input.URL()
	this.Data["EDPTenants"] = edpTenants
	this.TplName = "edp_tenants.tpl"
}

func (this *EDPTenantController) GetEDPComponents() {
	edpTenantName := this.GetString(":name")
	components := this.EDPTenantService.GetEDPComponents(edpTenantName)
	version, err := this.EDPTenantService.GetEDPVersionByName(edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["EDPTenantName"] = edpTenantName
	this.Data["EDPVersion"] = version
	this.Data["EDPComponents"] = components
	this.TplName = "edp_components.tpl"
}
