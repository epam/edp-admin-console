package config

import (
	"context"
	"time"

	"github.com/astaxie/beego"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const AuthEnable = "keycloakAuthEnabled"

type AppConfig struct {
	RunMode    string
	EDPVersion string
	AuthEnable bool
	BasePath   string
}

type KeycloakConfig struct {
	KeycloakURL      string
	ClientId         string
	ClientSecret     string
	Host             string
	CallBackEndpoint string
	StateAuthKey     string
}

type AuthController struct {
	Oauth2Service  oauth2.Config
	Verifier       *oidc.IDTokenVerifier
	StateAuthKey   string
	AuthSessionTTL time.Duration
	SessionTTL     time.Duration
}

//TODO remake config

func SetupConfig(_ context.Context, _ string) (*AppConfig, error) {
	authEnable, err := beego.AppConfig.Bool(AuthEnable)
	if err != nil {
		return nil, err
	}
	config := &AppConfig{
		RunMode:    beego.AppConfig.String("runmode"),
		EDPVersion: beego.AppConfig.String("edpVersion"),
		BasePath:   beego.AppConfig.String("basePath"),
		AuthEnable: authEnable,
	}
	return config, nil
}

func SetupAuthController(ctx context.Context, _ string) (*AuthController, error) {
	keycloakConfig, err := setupKeycloakConfig()
	if err != nil {
		return nil, err
	}
	provider, err := oidc.NewProvider(ctx, keycloakConfig.KeycloakURL)
	if err != nil {
		return nil, err
	}
	oauth2Config := oauth2.Config{
		ClientID:     keycloakConfig.ClientId,
		ClientSecret: keycloakConfig.ClientSecret,
		RedirectURL:  keycloakConfig.Host + keycloakConfig.CallBackEndpoint,
		Endpoint:     provider.Endpoint(),
	}
	oidcConfig := &oidc.Config{
		ClientID: keycloakConfig.ClientId,
	}
	authExpirationTime, err := beego.AppConfig.Int("authSessionTTLMinute")
	if err != nil {
		return nil, err
	}
	sessionTTL, err := beego.AppConfig.Int("sessionTTLMinute")
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(oidcConfig)
	authConfig := &AuthController{Oauth2Service: oauth2Config, Verifier: verifier,
		StateAuthKey: keycloakConfig.StateAuthKey, AuthSessionTTL: time.Duration(authExpirationTime) * time.Minute,
		SessionTTL: time.Duration(sessionTTL) * time.Minute}
	return authConfig, nil
}

func setupKeycloakConfig() (*KeycloakConfig, error) {
	config := &KeycloakConfig{
		KeycloakURL:      beego.AppConfig.String("keycloakURL"),
		ClientId:         beego.AppConfig.String("clientId"),
		ClientSecret:     beego.AppConfig.String("clientSecret"),
		Host:             beego.AppConfig.String("host"),
		CallBackEndpoint: beego.AppConfig.String("callBackEndpointV2"),
		StateAuthKey:     beego.AppConfig.String("stateAuthKey"),
	}
	return config, nil
}
