package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Error500() {
	this.TplName = "error_500.html"
}

func (this *ErrorController) Error403() {
	this.TplName = "error_403.html"
}
