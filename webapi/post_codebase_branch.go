package webapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"edp-admin-console/internal/applog"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

const paramWaitingForBranch = "waitingforbranch"

func (h HandlerEnv) CreateBranch(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	logger.Debug("in handler")

	codebaseName := chi.URLParam(request, "codebaseName")
	codebaseCR, errGet := h.NamespacedClient.GetCodebase(ctx, codebaseName)
	if errGet != nil {
		logger.Error("cant get codebase with that name", zap.Error(errGet))
		InternalErrorResponse(ctx, writer, "cant get codebase with that name")
		return
	}

	err := request.ParseForm()
	if err != nil {
		logger.Error("cant parse form", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant parse form")
		return
	}

	codebaseBranchName := request.Form.Get("name")
	commit := request.Form.Get("commit")
	version := request.Form.Get("version")
	postfix := request.Form.Get("versioningPostfix")

	release := request.Form.Get("releaseBranch")
	isRelease := false
	if release == "true" {
		isRelease = true
	}

	branchSpec := codeBaseApi.CodebaseBranchSpec{
		CodebaseName:     codebaseName,
		BranchName:       codebaseBranchName,
		FromCommit:       commit,
		Version:          GetVersionOrNil(version, postfix),
		Release:          isRelease,
		ReleaseJobParams: nil,
	}

	codebaseBranchCRName := fmt.Sprintf("%s-%s", codebaseName, strings.ReplaceAll(codebaseBranchName, "/", "-"))
	errBranch := h.NamespacedClient.CreateCBBranchByCustomFields(ctx, codebaseBranchCRName, branchSpec, codeBaseApi.CodebaseBranchStatus{})
	if errBranch != nil && !k8serrors.IsAlreadyExists(errBranch) {
		logger.Error("cant create codebase branch", zap.Error(errBranch))
		InternalErrorResponse(ctx, writer, "cant create codebase branch")
		return
	}

	if isRelease {
		masterVersion := request.Form.Get("masterVersion")
		snapshotStaticField := request.Form.Get("snapshotStaticField")
		defaultBranch := GetVersionOrNil(masterVersion, snapshotStaticField)
		if defaultBranch != nil {
			codebaseCR.Spec.DefaultBranch = *defaultBranch
		}
		errUpdate := h.NamespacedClient.UpdateCodebaseByCustomFields(ctx, codebaseName, codebaseCR.Spec, codebaseCR.Status)
		if errUpdate != nil {
			logger.Error("cant get codebase with that name", zap.Error(errGet))
			InternalErrorResponse(ctx, writer, "cant get codebase with that name")
			return
		}
	}

	if k8serrors.IsAlreadyExists(errBranch) {
		redirectUrl := fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?errorExistingBranch=%s#branchExistsModal", h.Config.BasePath,
			codebaseName, url.PathEscape(codebaseBranchName))
		http.Redirect(writer, request, redirectUrl, http.StatusFound)
		return
	}

	redirectUrl := fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?%s=%s#branchSuccessModal", h.Config.BasePath,
		codebaseName, paramWaitingForBranch, url.PathEscape(codebaseBranchName))
	http.Redirect(writer, request, redirectUrl, http.StatusFound)
}

func GetVersionOrNil(value, postfix string) *string {
	if value == "" {
		return nil
	}

	if postfix == "" {
		v := value
		return &v
	}

	v := fmt.Sprintf("%v-%v", value, postfix)

	return &v
}
