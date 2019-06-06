package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
)

type ServiceController struct {
	beego.Controller
	CatalogService service.CatalogService
}

func (s *ServiceController) GetServicePage() {
	services, err := s.CatalogService.GetAllServices()
	if err != nil {
		s.Abort("500")
		return
	}

	s.Data["EDPVersion"] = context.EDPVersion
	s.Data["Username"] = s.Ctx.Input.Session("username")
	s.Data["Services"] = services
	s.TplName = "service.html"
}
