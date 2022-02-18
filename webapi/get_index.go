package webapi

import (
	"net/http"
	"path"

	edpcontext "edp-admin-console/context"
)

func (h *HandlerEnv) Index(writer http.ResponseWriter, request *http.Request) {

	type data struct {
		BasePath   string
		EDPVersion string
		Username   string
	}

	tmplData := data{
		BasePath:   edpcontext.BasePath,
		EDPVersion: edpcontext.EDPVersion,
		Username:   "testname",
	}

	templatePaths := []string{
		path.Join(h.WorkingDir, "/viewsV2/index.html"),
		path.Join(h.WorkingDir, "/views/template/footer_template.html"),
		path.Join(h.WorkingDir, "/views/template/header_template.html"),
	}

	template := &Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "index.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(request.Context(), writer, template)
}
