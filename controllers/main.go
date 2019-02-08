package controllers

import "github.com/astaxie/beego"

type MainController struct {
	beego.Controller
}

func (this *MainController) Index() {
	this.TplName = "index.html"
}

func (this *MainController) Applications() {
	this.TplName = "applications.html"
}

func (this *MainController) AddingApplication() {
	this.TplName = "addingApplication.html"
}
