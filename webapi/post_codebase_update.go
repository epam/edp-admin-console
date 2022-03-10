package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"edp-admin-console/context"
	"edp-admin-console/internal/applog"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
)

func (h *HandlerEnv) PostCodebaseUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		logger.Error("parse form to update codebase failed", zap.Error(err))
		redirectURI := fmt.Sprintf("%v/v2/admin/edp/overview#codebaseUpdateErrorModal",
			context.BasePath)
		FoundRedirectResponse(w, r, redirectURI)
		return
	}

	codebaseName := r.Form.Get("name")
	codebaseCR, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
	if err != nil {
		logger.Error("get codebase in k8s failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		redirectURI := fmt.Sprintf("%v/v2/admin/edp/%v/overview#codebaseUpdateErrorModal",
			context.BasePath, codebaseTypeToMenuItem(codebaseCR.Spec.Type))
		FoundRedirectResponse(w, r, redirectURI)
		return
	}

	jiraServerToggle := r.Form.Get("jiraServerToggle")
	if jiraServerToggle == "on" {
		formData := r.Form
		jiraFieldNames := make([]string, 0)
		jiraPatterns := make([]string, 0)
		for k, v := range formData {
			if k == "jiraFieldName" {
				jiraFieldNames = append(jiraFieldNames, v...)
			}
			if k == "jiraPattern" {
				jiraPatterns = append(jiraPatterns, v...)
			}
		}
		payload, mappingErr := advanceMappingFormToJSON(jiraFieldNames, jiraPatterns)
		if mappingErr != nil {
			logger.Error("cant parse form to update codebase", zap.Error(mappingErr))
			redirectURI := fmt.Sprintf("%v/v2/admin/edp/%v/overview#codebaseUpdateErrorModal",
				context.BasePath, codebaseTypeToMenuItem(codebaseCR.Spec.Type))
			FoundRedirectResponse(w, r, redirectURI)
			return
		}
		commitMessagePattern := r.Form.Get("commitMessagePattern")
		ticketNamePattern := r.Form.Get("ticketNamePattern")
		jiraServer := r.Form.Get("jiraServer")
		codebaseCR.Spec.JiraIssueMetadataPayload = payload
		codebaseCR.Spec.CommitMessagePattern = strToPtr(commitMessagePattern)
		codebaseCR.Spec.TicketNamePattern = strToPtr(ticketNamePattern)
		codebaseCR.Spec.JiraServer = strToPtr(jiraServer)
	} else {
		codebaseCR.Name = codebaseName
	}

	err = h.NamespacedClient.Update(ctx, codebaseCR)
	if err != nil {
		logger.Error("update codebase in k8s failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		redirectURI := fmt.Sprintf("%v/v2/admin/edp/%v/overview#codebaseUpdateErrorModal",
			context.BasePath, codebaseTypeToMenuItem(codebaseCR.Spec.Type))
		FoundRedirectResponse(w, r, redirectURI)
		return
	}

	redirectURI := fmt.Sprintf("%v/v2/admin/edp/%v/overview#codebaseUpdateSuccessModal",
		context.BasePath, codebaseTypeToMenuItem(codebaseCR.Spec.Type))
	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func advanceMappingFormToJSON(jiraFieldNames, jiraPatterns []string) (*string, error) {
	if jiraFieldNames == nil && jiraPatterns == nil {
		return nil, nil
	}

	payload := make(map[string]string, len(jiraFieldNames))
	for i, name := range jiraFieldNames {
		payload[name] = jiraPatterns[i]
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return util.GetStringP(string(jsonPayload)), nil
}

func codebaseTypeToMenuItem(codebaseType string) string {
	if codebaseType == consts.Autotest {
		return "autotest"
	}
	return codebaseType
}
