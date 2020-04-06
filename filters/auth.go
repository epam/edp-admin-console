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

package filters

import (
	ctx "context"
	appCtx "edp-admin-console/context"
	"edp-admin-console/service/logger"
	"encoding/json"
	bgCtx "github.com/astaxie/beego/context"
	"github.com/coreos/go-oidc"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

var log = logger.GetLogger()

func AuthFilter(context *bgCtx.Context) {
	log.Debug("Start auth filter..")
	tsRaw := context.Input.Session("token_source")
	if tsRaw == nil {
		log.Debug("There are no token source in the session")
		startAuth(context)
		return
	}
	ts := tsRaw.(oauth2.TokenSource)
	token, err := ts.Token()
	if err != nil {
		log.Debug("Token source presented in the session is not valid")
		startAuth(context)
		return
	}
	idToken, err := appCtx.GetAuthConfig().Verifier.Verify(ctx.Background(), token.AccessToken)
	if err != nil {
		log.Debug("Token presented in the session is not valid")
		startAuth(context)
		return
	}
	realmRoles := getRealmRoles(context, idToken)
	log.Info("Roles have been retrieved from the token", zap.Strings("roles", realmRoles))
	context.Output.Session("realm_roles", realmRoles)
	username := getUserInfoFromToken(context, idToken, "name")
	log.Info("Username has been fetched from token", zap.String("username", username))
	context.Output.Session("username", username)
}

func getRealmRoles(context *bgCtx.Context, token *oidc.IDToken) []string {
	log.Debug("Start to check roles ...")
	var claim map[string]*json.RawMessage
	err := token.Claims(&claim)
	if err != nil {
		log.Error("Error has been occurred during the parsing token", zap.Any("token", token))
		context.Abort(200, "500")
	}
	var realmAccess map[string]*[]string
	err = json.Unmarshal(*claim["realm_access"], &realmAccess)
	if err != nil {
		log.Error("Error has been occurred during the parsing token", zap.Any("token", token))
		context.Abort(200, "500")
	}

	return *realmAccess["roles"]
}

func startAuth(context *bgCtx.Context) {
	authConfig := appCtx.GetAuthConfig()
	state := uuid.NewV4().String()
	log.Info("State has been generated, saved in the session and added in the auth request",
		zap.String("state", state))
	context.Output.Session(authConfig.StateAuthKey, state)
	if context.Request.Method == "GET" {
		context.Output.Session("request_path", context.Request.URL.Path)
	}
	context.Redirect(http.StatusFound, authConfig.Oauth2Config.AuthCodeURL(state))
}

func getUserInfoFromToken(context *bgCtx.Context, token *oidc.IDToken, userKey string) string {
	var claim map[string]*json.RawMessage
	err := token.Claims(&claim)
	if err != nil {
		log.Error("Error has been occurred during the parsing token", zap.Any("token", token))
		context.Abort(200, "500")
	}
	return strings.Replace(string(*claim[userKey]), "\"", "", -1)
}
