package context

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
)

type KeycloakParameters struct {
	KeycloakURL      string
	ClientId         string
	ClientSecret     string
	Host             string
	CallBackEndpoint string
	StateAuthKey     string
}

type AuthConfig struct {
	Oauth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	StateAuthKey string
}

var authConfig AuthConfig

func InitAuth() {
	parameters := getParameters()
	log.Println(fmt.Sprintf("Keycloak has been retrieved: %s. \n"+
		"ClientId has been retrieved: %s. \n"+
		"ClientSecret has been retrieved: %s. \n"+
		"Host has been retrieved: %s. \n"+
		"CallBackEndpoint has been retrieved: %s. \n"+
		"StateAuthKey has been retrieved: %s.",
		parameters.KeycloakURL, parameters.ClientId, parameters.ClientSecret,
		parameters.Host, parameters.CallBackEndpoint, parameters.StateAuthKey))

	provider, err := oidc.NewProvider(context.Background(), parameters.KeycloakURL)

	if err != nil {
		log.Fatal(err)
		return
	}

	oauth2Config := oauth2.Config{
		ClientID:     parameters.ClientId,
		ClientSecret: parameters.ClientSecret,
		RedirectURL:  parameters.Host + parameters.CallBackEndpoint,
		Endpoint:     provider.Endpoint(),
	}

	oidcConfig := &oidc.Config{
		ClientID: parameters.ClientId,
	}

	verifier := provider.Verifier(oidcConfig)

	authConfig = AuthConfig{oauth2Config, verifier, parameters.StateAuthKey}
}

func getParameters() KeycloakParameters {
	//todo add checking all variables
	return KeycloakParameters{
		beego.AppConfig.String("keycloakURL"),
		beego.AppConfig.String("clientId"),
		beego.AppConfig.String("clientSecret"),
		beego.AppConfig.String("host"),
		beego.AppConfig.String("callBackEndpoint"),
		beego.AppConfig.String("stateAuthKey")}
}

func GetAuthConfig() AuthConfig {
	return authConfig
}
