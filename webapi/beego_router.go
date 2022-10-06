/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package webapi

import (
	"context"
	"fmt"
	"path"

	"github.com/astaxie/beego"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	edpcontext "edp-admin-console/context"
	"edp-admin-console/filters"
	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	"edp-admin-console/service/logger"
)

var zaplog = logger.GetLogger()

const (
	apiV2Scope = "/api/v2"
	edpScope   = "/edp"
)

func SetupRouter(namespacedClient *k8s.RuntimeNamespacedClient, workingDir string, confV2 *config.AppConfig, clusterConfig *rest.Config) {
	zaplog.Info(
		"Start application...",
		zap.String("mode", confV2.RunMode),
		zap.String("edp version", confV2.EDPVersion),
	)

	permissions := filters.PermissionsMap()
	accessHandlerEnv := &filters.AccessControlEnv{
		Permissions: permissions,
	}

	if confV2.AuthEnable {
		edpcontext.InitAuth()

		// auth and role access for v2 api
		beego.InsertFilter(fmt.Sprintf("%s%s%s/*", edpcontext.BasePath, apiV2Scope, edpScope), beego.BeforeRouter, filters.AuthRestFilter)
		beego.InsertFilter(fmt.Sprintf("%s%s%s/*", edpcontext.BasePath, apiV2Scope, edpScope), beego.BeforeRouter, accessHandlerEnv.RoleAccessControlRestFilter)
	} else {
		beego.InsertFilter(fmt.Sprintf("%s/*", edpcontext.BasePath), beego.BeforeRouter, filters.StubAuthFilter)
	}

	v2APIHandler := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(confV2), WithClusterConfig(clusterConfig))

	authOpts := make([]HandlerAuthOption, 0)
	authOpts = append(authOpts, WithBasePath(confV2.BasePath))
	if v2APIHandler.Config.AuthEnable {
		authController, err := config.SetupAuthController(context.Background(), "conf/app.conf")
		if err != nil {
			zaplog.Error("cant setup authController", zap.Error(err))
		}
		authOpts = append(authOpts, WithAuthController(authController))
	}
	v2APIAuthHandler := HandlerAuthWithOption(authOpts...)
	v2APIRouter := V2APIRouter(v2APIHandler, v2APIAuthHandler, zaplog)

	// see https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md#handler-register
	// and isPrefix parameter
	beego.Handler(path.Join(v2APIHandler.Config.BasePath, "/v2"), v2APIRouter, true)
	beego.Handler(path.Join(v2APIHandler.Config.BasePath, apiV2Scope, edpScope), v2APIRouter, true)
}

func V2APIRouter(handlerEnv *HandlerEnv, authHandler *HandlerAuth, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		WithLoggerMw(logger),
		WithLogRequestBoundaries(),
	)

	basePath := path.Join(handlerEnv.Config.BasePath, "/")
	router.Route(basePath, func(baseRouter chi.Router) {
		baseRouter.Route(apiV2Scope, func(v2APIRouter chi.Router) {
			v2APIRouter.Route(edpScope, func(edpScope chi.Router) {
				edpScope.Route("/cd-pipeline", func(pipelinesRoute chi.Router) {
					pipelinesRoute.Route("/{pipelineName}", func(pipelineRoute chi.Router) {
						pipelineRoute.Get("/", handlerEnv.GetPipeline)
						pipelineRoute.Get("/stage/{stageName}", handlerEnv.GetStagePipeline)
					})
				})
				edpScope.Route("/codebase", func(codebasesRoute chi.Router) {
					codebasesRoute.Get("/", handlerEnv.GetCodebases)
					codebasesRoute.Get("/{codebaseName}", handlerEnv.GetCodebase)
				})
			})
		})
	})

	return router
}
