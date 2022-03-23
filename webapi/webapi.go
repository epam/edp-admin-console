package webapi

import (
	"context"
	"html/template"

	"golang.org/x/oauth2"
	"k8s.io/client-go/rest"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
)

type HandlerEnv struct {
	NamespacedClient *k8s.RuntimeNamespacedClient
	FuncMap          template.FuncMap
	WorkingDir       string
	Config           *config.AppConfig
	ClusterConfig    *rest.Config
}

type HandlerAuth struct {
	StateMap       map[string]string
	TokenMap       map[string]oauth2.TokenSource
	UrlMap         map[string]string
	AuthController *config.AuthController
	BasePath       string
}

type HandlerAuthOption func(handler *HandlerAuth)

func WithBasePath(basePath string) HandlerAuthOption {
	return func(handler *HandlerAuth) {
		handler.BasePath = basePath
	}
}

func WithAuthController(controller *config.AuthController) HandlerAuthOption {
	return func(handler *HandlerAuth) {
		handler.AuthController = controller
	}
}

func HandlerAuthWithOption(opts ...HandlerAuthOption) *HandlerAuth {
	stateMap := make(map[string]string)
	tokenMap := make(map[string]oauth2.TokenSource)
	urlMap := make(map[string]string)
	handler := &HandlerAuth{
		StateMap: stateMap,
		TokenMap: tokenMap,
		UrlMap:   urlMap,
	}
	for i := range opts {
		opts[i](handler)
	}
	return handler
}

func CreateCommonFuncMap() template.FuncMap {
	return template.FuncMap{
		"getCurrentYear":          getCurrentYear,
		"add":                     add,
		"getDefaultBranchVersion": getDefaultBranchVersion,
		"incrementVersion":        incrementVersion,
		"compareJiraServer":       compareJiraServer,
		"params":                  params,
		"capitalizeFirst":         CapitalizeFirstLetter,
		"capitalizeAll":           CapitalizeAll,
		"lowercaseAll":            LowercaseAll,
	}
}

type HandlerEnvOption func(handler *HandlerEnv)

func WithConfig(config *config.AppConfig) HandlerEnvOption {
	return func(handler *HandlerEnv) {
		handler.Config = config
	}
}

func WithClient(client *k8s.RuntimeNamespacedClient) HandlerEnvOption {
	return func(handler *HandlerEnv) {
		handler.NamespacedClient = client
	}
}

func WithFuncMap(funcMap template.FuncMap) HandlerEnvOption {
	return func(handler *HandlerEnv) {
		handler.FuncMap = funcMap
	}
}

func WithClusterConfig(config *rest.Config) HandlerEnvOption {
	return func(handler *HandlerEnv) {
		handler.ClusterConfig = config
	}
}

func WithWorkingDir(workingDir string) HandlerEnvOption {
	return func(handler *HandlerEnv) {
		handler.WorkingDir = workingDir
	}
}

func NewHandlerEnv(opts ...HandlerEnvOption) *HandlerEnv {

	handler := &HandlerEnv{}
	for i := range opts {
		opts[i](handler)
	}
	return handler
}

func UserFromContext(ctx context.Context) AuthorisedUser {
	user, ok := ctx.Value(AuthorisedUserKey{}).(AuthorisedUser)
	if ok {
		return user
	}
	user = GuestUser()
	return user

}
