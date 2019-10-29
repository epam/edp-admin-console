/*
 * Copyright 2019 EPAM Systems.
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

package service

import (
	"edp-admin-console/context"
	"edp-admin-console/k8s"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
)

type EDPTenantService struct {
	Clients k8s.ClientSet
}

func (this EDPTenantService) GetVcsIntegrationValue() (bool, error) {
	coreClient := this.Clients.CoreClient

	res, err := coreClient.ConfigMaps(context.Namespace).Get("user-settings", metav1.GetOptions{})

	if err != nil {
		log.Printf("An error has occurred while getting user settings: %s", err)
		return false, err
	}

	var vcsEnabled = res.Data["vcs_integration_enabled"]

	if len(vcsEnabled) == 0 {
		log.Println("vcs_integration_enabled property doesn't exist")
		return false, errors.New("NOT_FOUND")
	}

	result, err := strconv.ParseBool(vcsEnabled)

	if err != nil {
		log.Printf("An error has occurred while parsing 'vcs_integration_enabled=false': %s", err)
		return false, err
	}
	return result, nil
}
