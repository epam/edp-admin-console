package controllers

import "github.com/astaxie/beego"

type MainController struct {
	beego.Controller
}

func (this *MainController) Index() {
	this.TplName = "index.html"
}

func (this *MainController) GetApplicationPage() {
	this.TplName = "application.html"
}

func (this *MainController) GetCreateApplicationPage() {
	this.TplName = "create_application.html"
}
