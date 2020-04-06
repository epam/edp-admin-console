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
	bgCtx "github.com/astaxie/beego/context"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strings"
)

func AuthRestFilter(context *bgCtx.Context) {
	log.Debug("Start auth rest filter..")
	token := context.Input.Header("Authorization")
	if token == "" {
		log.Error("There are no token in the session")
		http.Error(context.ResponseWriter, "The request header doesn't contain token.", http.StatusBadRequest)
		return
	}

	token, err := tryToRemoveBearerPrefix(token)
	if err != nil {
		log.Error("An error has occurred while checking regexp.")
		http.Error(context.ResponseWriter, "Internal Error", http.StatusInternalServerError)
		return
	}

	idToken, err := appCtx.GetAuthConfig().Verifier.Verify(ctx.Background(), token)
	if err != nil {
		log.Error("Token presented in the session is not valid")
		http.Error(context.ResponseWriter, "Token presented in the session is not valid", http.StatusUnauthorized)
		return
	}

	realmRoles := getRealmRoles(context, idToken)
	log.Info("Roles have been retrieved from the token", zap.Strings("roles", realmRoles))
	usr := getUserInfoFromToken(context, idToken, "preferred_username")
	log.Info("Username has been fetched from token", zap.String("username", usr))
	context.Output.Session("realm_roles", realmRoles)
	context.Output.Session("username", usr)
}

func tryToRemoveBearerPrefix(token string) (string, error) {
	isMatched, err := regexp.MatchString("^Bearer", token)
	if err != nil {
		return "", err
	}

	if isMatched {
		return strings.Replace(token, "Bearer", "", -1), nil
	}
	return token, nil
}
