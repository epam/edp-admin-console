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
	"edp-admin-console/k8s"
	"edp-admin-console/models"
	"edp-admin-console/repository"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
	"strings"
)

type EDPTenantService struct {
	IEDPTenantRep repository.IEDPTenantRepository
	Clients       k8s.ClientSet
}

var (
	edpComponentNames = []string{"Jenkins", "Gerrit", "Sonar", "Nexus", "EDP-Cockpit"}
	wildcard          = beego.AppConfig.String("dnsWildcard")
)

func (this EDPTenantService) GetEDPTenants(resourceAccess map[string][]string) ([]*models.EDPTenant, error) {
	edpTenantNames := filterEdpTenantNamesWithoutSuffixWithCurrentRoles(resourceAccess)
	if edpTenantNames == nil {
		log.Println("There aren't edp tenants to display.")
		return nil, nil
	}

	edpSpecs, err := this.IEDPTenantRep.GetAllEDPTenantsByNames(edpTenantNames)
	if err != nil {
		log.Printf("Couldn't get all EDP specifications. Reason: %v\n", err)
		return nil, err
	}

	return edpSpecs, nil
}

func (this EDPTenantService) GetEDPVersionByName(edpTenantName string) (string, error) {
	version, err := this.IEDPTenantRep.GetEdpVersionByName(edpTenantName)
	if err != nil {
		log.Printf("An error has occurred while getting version of %s EDP.", edpTenantName)
		return "", err
	}
	return version, nil
}

func (this EDPTenantService) GetTenantByName(edpName string) (*models.EDPTenant, error) {
	edpTenant, err := this.IEDPTenantRep.GetTenantByName(edpName)
	if err != nil {
		log.Printf("An error has occurred while getting tenant by %s name.", edpName)
		return nil, err
	}
	return edpTenant, nil
}

func (edpService EDPTenantService) GetEDPComponents(edpTenantName string) map[string]string {
	var compWithLinks = make(map[string]string, len(edpComponentNames))
	for _, val := range edpComponentNames {
		compWithLinks[val] = fmt.Sprintf("https://%s-%s-edp-cicd.%s", strings.ToLower(val), edpTenantName, wildcard)
	}
	return compWithLinks
}

func (this EDPTenantService) GetVcsIntegrationValue(edpName string) (bool, error) {
	coreClient := this.Clients.CoreClient
	namespace := edpName + "-edp-cicd"

	res, err := coreClient.ConfigMaps(namespace).Get("user-settings", metav1.GetOptions{})

	if err != nil {
		log.Printf("An error has occurred while getting user settings: %s", err)
		return false, err
	}

	var vcsEnabled = res.Data["vcs_integration_enabled"]

	if len(vcsEnabled) == 0 {
		log.Println("vcs_integration_enabled property doesn't exist. Configured default value: 'vcs_integration_enabled=false'")
		return false, nil
	}

	result, err := strconv.ParseBool(vcsEnabled)

	if err != nil {
		log.Printf("An error has occurred while parsing 'vcs_integration_enabled=false': %s", err)
		return false, err
	}

	return result, nil
}

func filterEdpTenantNamesWithoutSuffixWithCurrentRoles(resourceAccess map[string][]string) []string {
	var edpTenants []string
	suffix := "-edp"
	for key, value := range resourceAccess {
		if strings.HasSuffix(key, suffix) &&
			(util.Contains(value, beego.AppConfig.String("adminRole")) || util.Contains(value, beego.AppConfig.String("developerRole"))) {
			edpTenants = append(edpTenants, strings.TrimSuffix(key, suffix))
		}
	}
	return edpTenants
}
