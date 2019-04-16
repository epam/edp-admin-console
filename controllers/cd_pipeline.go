package controllers

import (
	"edp-admin-console/models"
	"github.com/astaxie/beego"
)

type CDPipelineRestController struct {
	beego.Controller
}

func (this *CDPipelineRestController) GetCDPipelineByName() {
	cdPipeline := models.CDPipelineReadRestApi{
		Name: "stub",
		CodebaseBranches: []models.CodebaseBranchReadRestApi{
			{
				BranchName: "stub-branch-release-1.0",
				AppName:    "stub-app-1",
			},
			{
				BranchName: "stub-branch-release-2.0",
				AppName:    "stub-app-2",
			},
		},
	}
	this.Data["json"] = cdPipeline
	this.ServeJSON()
}
