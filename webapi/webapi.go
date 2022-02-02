package webapi

import (
	"context"
	"log"

	"go.uber.org/zap"

	"edp-admin-console/k8s"
)

type HandlerEnv struct {
	NamespacedClient *k8s.RuntimeNamespacedClient
}

func NewHandlerEnv(namespacedClient *k8s.RuntimeNamespacedClient) *HandlerEnv {
	return &HandlerEnv{NamespacedClient: namespacedClient}
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
