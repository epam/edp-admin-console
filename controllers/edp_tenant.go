package controllers

import (
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
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
	this.TplName = "edp_tenants.html"
}

func (this *EDPTenantController) GetEDPComponents() {
	resourceAccess := this.Ctx.Input.Session("resource_access").(map[string][]string)
	edpTenants, err := this.EDPTenantService.GetEDPTenants(resourceAccess)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	edpTenantName := this.GetString(":name")
	components := this.EDPTenantService.GetEDPComponents(edpTenantName)
	version, err := this.EDPTenantService.GetEDPVersionByName(edpTenantName)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["InputURL"] = strings.TrimSuffix(this.Ctx.Input.URL(), "/"+edpTenantName)
	this.Data["EDPTenantName"] = edpTenantName
	this.Data["EDPVersion"] = version
	this.Data["EDPComponents"] = components
	this.Data["EDPTenants"] = edpTenants
	this.TplName = "edp_components.html"
}

func (this *EDPTenantController) GetVcsIntegrationValue() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue()

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = isVcsEnabled
	this.ServeJSON()
}