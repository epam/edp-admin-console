package webapi

import (
	"fmt"
	"net/http"

	"edp-admin-console/context"
	"edp-admin-console/internal/applications"
	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/imagestream"

	"go.uber.org/zap"
)

func (h *HandlerEnv) PostCodebase(w http.ResponseWriter, r *http.Request) { // it is delete request, really. sorry
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		logger.Error("parse form failed", zap.Error(err))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/overview#codebaseDeleteErrorModal", context.BasePath),
			http.StatusFound)
		return
	}
	codebaseNameToDelete := r.Form.Get("name")
	codebaseType := r.Form.Get("codebase-type")
	k8sClient := h.NamespacedClient

	codebaseCR, err := applications.ByNameIFExists(ctx, k8sClient, codebaseNameToDelete)
	if err != nil {
		logger.Error("get codebase failed", zap.Error(err), zap.String("codebase_name", codebaseNameToDelete))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v#codebaseDeleteErrorModal",
				context.BasePath, codebaseType, codebaseNameToDelete),
			http.StatusFound)
		return
	}
	if codebaseCR == nil {
		logger.Error("codebase not found", zap.String("codebase_name", codebaseNameToDelete))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview#codebaseDeleteErrorModal", context.BasePath, codebaseType),
			http.StatusFound)
		return
	}

	codebaseStreamsMap, err := imagestream.OutputCBStreamsForCodebaseNames(ctx, k8sClient, []string{codebaseNameToDelete})
	if err != nil {
		logger.Error("get codebase image streams failed", zap.Error(err), zap.String("codebase_name", codebaseNameToDelete))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v#codebaseDeleteErrorModal",
				context.BasePath, codebaseType, codebaseNameToDelete),
			http.StatusFound)
		return
	}

	cdPipelinesCR, err := h.NamespacedClient.CDPipelineList(ctx)
	if err != nil {
		logger.Error("get cd pipelines failed", zap.Error(err), zap.String("codebase_name", codebaseNameToDelete))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v#codebaseDeleteErrorModal",
				context.BasePath, codebaseType, codebaseNameToDelete),
			http.StatusFound)
		return
	}

	codebaseIsInUse := false
	errmsg := ""
	for cbName, cbImageStreams := range codebaseStreamsMap {
		for _, cbIS := range cbImageStreams {
			for _, cdPipeline := range cdPipelinesCR {
				for _, cdPipelineIS := range cdPipeline.Spec.InputDockerStreams { // this is potentially heavy operation and should be optimized
					if cdPipelineIS == cbIS.Name {
						logger.Info("codebase image stream is used by cd pipeline",
							zap.String("codebase_image_stream", cbIS.Name),
							zap.String("cd_pipeline", cdPipeline.Name),
							zap.String("codebase_name", cbName),
						)
						codebaseIsInUse = true
						errmsg = fmt.Sprintf(
							"application '%s' and its image stream '%s' are used by '%s' pipeline",
							codebaseNameToDelete, cdPipelineIS, cdPipeline.Name)
					}
				}
			}
		}
	}

	if codebaseIsInUse {
		logger.Error("codebase is in use", zap.String("codebase_name", codebaseNameToDelete))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v&errmsg=%s#codebaseIsUsed",
				context.BasePath, codebaseType, codebaseNameToDelete, errmsg),
			http.StatusFound)
		return
	}

	deleteErr := k8sClient.DeleteCodebase(ctx, codebaseCR)
	if err != nil {
		logger.Error("delete codebase custom resource failed", zap.Error(deleteErr), zap.String("codebase_name", codebaseCR.Name))
		http.Redirect(w, r,
			fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v#codebaseDeleteErrorModal",
				context.BasePath, codebaseType, codebaseNameToDelete),
			http.StatusFound)
		return
	}

	http.Redirect(w, r,
		fmt.Sprintf("%s/v2/admin/edp/%v/overview?codebase=%v#codebaseIsDeleted", context.BasePath, codebaseType, codebaseNameToDelete),
		http.StatusFound)
}
