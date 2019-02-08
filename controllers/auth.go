package controllers

import (
	"context"
	ctx "edp-admin-console/context"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

type AuthController struct {
	beego.Controller
}

func (this *AuthController) Callback() {
	authConfig := ctx.GetAuthConfig()
	log.Println("Start callback flow...")
	queryState := this.Ctx.Input.Query("state")
	log.Printf("State %s has been retrived from query param", queryState)
	sessionState := this.Ctx.Input.Session(authConfig.StateAuthKey)
	log.Printf("State %s has been retrived from the session", sessionState)
	if queryState != sessionState {
		http.Error(this.Ctx.ResponseWriter, "State does not match", http.StatusBadRequest)
		return
	}

	authCode := this.Ctx.Input.Query("code")
	log.Println("Authorization code has been retrieved from query param")
	token, err := authConfig.Oauth2Config.Exchange(context.Background(), authCode)

	if err != nil {
		log.Printf("Failed to exchange token with code: %s", authCode)
		http.Error(this.Ctx.ResponseWriter, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Authorization code has been successfully exchanged with token")

	this.Ctx.Output.Session("token", token)
	log.Println("Token has been saved to the session")
	this.Redirect("/admin/edp", 302)
}
