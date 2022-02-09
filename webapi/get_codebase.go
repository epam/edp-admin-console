package webapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type GetCodebaseResponse GetCodebase

func (h *HandlerEnv) GetCodebase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)
	codebaseName := chi.URLParam(r, "codebaseName")
	if codebaseName == "" {
		logger.Error("codebase not passed")
		BadRequestResponse(ctx, w, "codebase not passed")
		return
	}

	crCodebase, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
	if err != nil {
		logger.Error("get codebase by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		InternalErrorResponse(ctx, w, "get codebase by name failed")
		return
	}
	codebaseResponse := GetCodebase{
		Name:         crCodebase.Name,
		GitServer:    crCodebase.Spec.GitServer,
		BuildTool:    strings.ToLower(crCodebase.Spec.BuildTool),
		Type:         crCodebase.Spec.Type,
		EmptyProject: crCodebase.Spec.EmptyProject,
		JenkinsSlave: pointerToStr(crCodebase.Spec.JenkinsSlave),
	}

	response := GetCodebaseResponse(codebaseResponse)
	OKJsonResponse(ctx, w, response)
}

func pointerToStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
