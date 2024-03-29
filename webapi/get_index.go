package webapi

import (
	"net/http"
	"path"
)

func (h *HandlerEnv) Index(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	user := UserFromContext(ctx)

	type data struct {
		BasePath   string
		EDPVersion string
		Username   string
	}

	tmplData := data{
		BasePath:   h.Config.BasePath,
		EDPVersion: h.Config.EDPVersion,
		Username:   user.UserName(),
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
