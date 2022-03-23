package webapi

import (
	"fmt"
	"net/http"
	"strings"

	"edp-admin-console/internal/applications"
	"edp-admin-console/internal/applog"

	"go.uber.org/zap"
)

func (h *HandlerEnv) DeleteCodebaseBranch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		logger.Error("parse form failed", zap.Error(err))
		InternalErrorResponse(ctx, w, "parse form failed")
		return
	}

	codebaseName := r.FormValue("codebase-name")
	cbBranchName := r.FormValue("name")
	cbBranchCRName := buildCobaseBranchCRName(codebaseName, cbBranchName)
	cbBranchCR, err := applications.CodebaseBranchByNameIfExists(ctx, h.NamespacedClient, cbBranchCRName)
	if err != nil {
		logger.Error("get codebase branch failed", zap.Error(err), zap.String("codebase_branch", cbBranchCRName))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?name=%s#branchDeletionErrorModal", h.Config.BasePath, codebaseName, cbBranchName),
			http.StatusFound)
		return
	}
	if cbBranchCR == nil {
		logger.Error("codebase branch not found", zap.String("codebase_branch", cbBranchCRName))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?name=%s#branchDeletionErrorModal", h.Config.BasePath, codebaseName, cbBranchName),
			http.StatusFound)
		return
	}

	/*
		Attention!!!
		Codebase branch can be deleted only if it is not used in any cd pipeline.
		At the moment there is no link "codebase branch - cd pipeline" in Custom Resources
		Please, fix this link ASAP and fix this handler.
	*/

	err = h.NamespacedClient.DeleteCodebaseBranch(ctx, cbBranchCR)
	if err != nil {
		logger.Error("delete codebase branch failed", zap.Error(err), zap.String("codebase_branch", cbBranchCRName))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?name=%s#branchDeletionErrorModal", h.Config.BasePath, codebaseName, cbBranchName),
			http.StatusFound)
		return
	}

	http.Redirect(w, r,
		fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?name=%s#branchDeletedSuccessModal", h.Config.BasePath, codebaseName, cbBranchName),
		http.StatusFound)
}

func buildCobaseBranchCRName(codebaseName, codebaseBranchName string) string {
	return fmt.Sprintf("%s-%s", codebaseName, strings.ReplaceAll(codebaseBranchName, "/", "-"))
}
