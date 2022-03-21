package webapi

import (
	"net/http"
	"path"
	"strings"

	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
)

type edpComponent struct {
	Visible bool
	Url     string
	Type    string
	Icon    string
}

type overviewData struct {
	EDPVersion         string
	Type               string
	BasePath           string
	DiagramPageEnabled bool
	EDPTenantName      string
	EDPComponents      []edpComponent
	Username           string
	InputURL           string
}

const HttpsScheme = "https://"
const HttpScheme = "http://"

func (h *HandlerEnv) GetOverviewPage(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)

	edpComponentCRs, err := h.NamespacedClient.EDPComponentList(ctx)
	if err != nil {
		logger.Error("cant get edp components list", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get edp components list")
		return
	}
	var edpComponents []edpComponent
	for i := range edpComponentCRs {
		url := edpComponentCRs[i].Spec.Url
		if edpComponentCRs[i].Spec.Type == consts.Openshift || edpComponentCRs[i].Spec.Type == consts.Kubernetes {
			url = util.CreateNativeProjectLink(url, h.NamespacedClient.Namespace)
		}
		if !strings.HasPrefix(url, HttpScheme) && !strings.HasPrefix(url, HttpsScheme) {
			url = HttpsScheme + url
		}
		edpComp := edpComponent{
			Visible: edpComponentCRs[i].Spec.Visible,
			Url:     url,
			Type:    edpComponentCRs[i].Spec.Type,
			Icon:    edpComponentCRs[i].Spec.Icon,
		}
		edpComponents = append(edpComponents, edpComp)
	}

	tmplData := overviewData{
		EDPVersion:         h.Config.EDPVersion,
		Type:               consts.Overview,
		BasePath:           h.Config.BasePath,
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		EDPTenantName:      h.NamespacedClient.Namespace,
		EDPComponents:      edpComponents,
		Username:           user.UserName(),
		InputURL:           request.URL.Path,
	}

	templatePaths := []string{
		path.Join(h.WorkingDir, "/viewsV2/edp_components.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/footer_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/header_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/navbar_template.html"),
	}

	template := &Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "edp_components.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, writer, template)

}
