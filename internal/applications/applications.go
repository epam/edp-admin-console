package applications

import (
	"context"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/k8s"
	"edp-admin-console/util/consts"
)

func ActiveApplications(cbList []codeBaseApi.Codebase) ([]codeBaseApi.Codebase, error) {
	return CodebasesByTypeAndStatus(cbList, consts.Application, consts.ActiveValue)
}

func ActiveAutotests(cbList []codeBaseApi.Codebase) ([]codeBaseApi.Codebase, error) {
	return CodebasesByTypeAndStatus(cbList, consts.Autotest, consts.ActiveValue)
}

func ActiveGroovyLibs(cbList []codeBaseApi.Codebase) ([]codeBaseApi.Codebase, error) {
	libs, err := CodebasesByTypeAndStatus(cbList, consts.Library, consts.ActiveValue)
	if err != nil {
		return nil, err
	}
	groovyLang := "groovy-pipeline"
	groovyLibs := make([]codeBaseApi.Codebase, 0)
	for _, lib := range libs {
		if lib.Spec.Lang == groovyLang {
			groovyLibs = append(groovyLibs, lib)
		}
	}
	return groovyLibs, nil
}

func CodebasesByTypeAndStatus(cbList []codeBaseApi.Codebase, cbType, cbStatus string) ([]codeBaseApi.Codebase, error) {
	activeApps := make([]codeBaseApi.Codebase, 0)
	for i := range cbList {
		codebaseCR := cbList[i]
		if codebaseCR.Spec.Type == cbType {
			if codebaseCR.Status.Value == cbStatus {
				activeApps = append(activeApps, codebaseCR)
			}
		}
	}
	return activeApps, nil
}

func ActiveCodebaseBranches(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, codebaseName string) ([]*codeBaseApi.CodebaseBranch, error) {
	logger := applog.LoggerFromContext(ctx)
	cbBranches, err := k8sClient.CodebaseBranchesListByCodebaseName(ctx, codebaseName)
	if err != nil {
		logger.Error("list codebase branches for codebase failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		return nil, err
	}

	activeBranches := make([]*codeBaseApi.CodebaseBranch, 0)
	for i := range cbBranches {
		branch := cbBranches[i]
		if branch.Status.Value == consts.ActiveValue {
			activeBranches = append(activeBranches, branch)
		}
	}
	return activeBranches, nil
}

func AppNameByInputIS(ctx context.Context, client *k8s.RuntimeNamespacedClient, isName string) (string, error) {
	stream, err := client.GetCodebaseImageStream(ctx, isName)
	if err != nil {
		return "", err
	}
	return stream.Spec.Codebase, nil
}
