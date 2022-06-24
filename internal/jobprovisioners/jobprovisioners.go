package jobprovisioners

import (
	"context"

	"go.uber.org/zap"

	jenkinsAPI "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1"

	"edp-admin-console/internal/applog"
	"edp-admin-console/k8s"
)

func ListNames(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient) ([]string, error) {
	jobProvisioners, err := List(ctx, k8sClient)
	if err != nil {
		return nil, err
	}
	namesMap := make(map[string]struct{})
	for _, jobProvision := range jobProvisioners {
		if _, ok := namesMap[jobProvision.Name]; !ok {
			namesMap[jobProvision.Name] = struct{}{}
		}
	}
	names := make([]string, 0)
	for k := range namesMap {
		names = append(names, k)
	}
	return names, nil
}

func List(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient) ([]jenkinsAPI.JobProvision, error) {
	logger := applog.LoggerFromContext(ctx)
	jenkinsList, err := k8sClient.GetJenkinsList(ctx)
	if err != nil {
		logger.Error("get jenkins list failed", zap.Error(err))
		return nil, err
	}

	jobProvisioners := make([]jenkinsAPI.JobProvision, 0)
	for _, jenkinsCR := range jenkinsList.Items {
		jobProvisioners = append(jobProvisioners, jenkinsCR.Status.JobProvisions...)
	}
	return jobProvisioners, nil
}
