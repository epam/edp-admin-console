package webapi

import (
	"edp-admin-console/k8s"
)

type HandlerEnv struct {
	NamespacedClient *k8s.RuntimeNamespacedClient
}

func NewHandlerEnv(namespacedClient *k8s.RuntimeNamespacedClient) *HandlerEnv {
	return &HandlerEnv{NamespacedClient: namespacedClient}
}
