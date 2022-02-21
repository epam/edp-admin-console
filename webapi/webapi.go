package webapi

import (
	"context"
	"html/template"
	"log"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
)

type HandlerEnv struct {
	NamespacedClient *k8s.RuntimeNamespacedClient
	FuncMap          template.FuncMap
	WorkingDir       string
	Config           *config.AppConfig
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

func getCurrentYear() int {
	return time.Now().Year()
}

func CreateCommonFuncMap() template.FuncMap {
	return template.FuncMap{
		"getCurrentYear": getCurrentYear,
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

type logCtx struct{}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logCtx{}, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	v, ok := ctx.Value(logCtx{}).(*zap.Logger)
	if !ok {
		log.Printf("logger not found: %+v (%T)", v, v)
		logger, err := zap.NewProduction(zap.WithCaller(true))
		if err != nil {
			log.Printf("init production logger failed: %s", err)
			return zap.NewExample() // fallback to simple logger
		}
		return logger
	}
	return v
}
