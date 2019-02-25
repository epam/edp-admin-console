package controllers

import (
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"net/http"
)

type ClusterController struct {
	beego.Controller
	ClusterService service.ClusterService
}

func (this *ClusterController) GetAllStorageClasses() {
	storageClasses, err := this.ClusterService.GetAllStorageClasses()

	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = storageClasses
	this.ServeJSON()
}
