package service

import (
	"edp-admin-console/models"
	"edp-admin-console/repository"
	"edp-admin-console/util"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"strings"
)

type EDPTenantService struct {
	EDPTenantRep repository.EDPTenantRep
}

var (
	edpComponentNames = []string{"Jenkins", "Gerrit", "Sonar", "Nexus", "Cockpit"}
	wildcard          = beego.AppConfig.String("wildcard")
)

func (edpService EDPTenantService) GetEDPTenants(resourceAccess map[string][]string) ([]*models.EDPTenant, error) {
	edpTenantNames := getEdpTenantNamesWithoutSuffix(resourceAccess)
	if edpTenantNames == nil {
		log.Println("There aren't edp tenants to display.")
		return nil, nil
	}

	edpSpecs, err := edpService.EDPTenantRep.GetAllEDPTenantsByNames(edpTenantNames)
	if err != nil {
		log.Printf("Couldn't get all EDP specifications. Reason: %v\n", err)
		return nil, err
	}

	return edpSpecs, nil
}

func (edpService EDPTenantService) GetEDPVersionByName(edpTenantName string) (string, error) {
	version, err := edpService.EDPTenantRep.GetEdpVersionByName(edpTenantName)
	if err != nil {
		log.Printf("An error has occurred while getting version of %s EDP.", edpTenantName)
		return "", err
	}
	return version, nil
}

func (edpService EDPTenantService) GetEDPComponents(edpTenantName string) map[string]string {
	var compWithLinks = make(map[string]string, len(edpComponentNames))
	for _, val := range edpComponentNames {
		link := fmt.Sprintf("https://%s-%s-edp-cicd.%s", strings.ToLower(val), edpTenantName, wildcard)
		compWithLinks[val] = link
	}
	return compWithLinks
}

func getEdpTenantNamesWithoutSuffix(resourceAccess map[string][]string) []string {
	var edpTenants []string
	suffix := "-edp"
	for key, value := range resourceAccess {
		if strings.HasSuffix(key, suffix) && util.Contains(value, beego.AppConfig.String("adminRole")) {
			edpTenants = append(edpTenants, strings.TrimSuffix(key, suffix))
		}
	}
	return edpTenants
}
