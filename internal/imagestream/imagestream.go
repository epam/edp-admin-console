package imagestream

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"go.uber.org/zap"

	"edp-admin-console/internal/applications"
	"edp-admin-console/internal/applog"
	"edp-admin-console/k8s"
)

const (
	verifiedImageStreamSuffix      = "verified"
	PreviousStageNameAnnotationKey = "deploy.edp.epam.com/previous-stage-name"
)

type EmptyImageStreamErr struct {
	crName string
}

func (err *EmptyImageStreamErr) Error() string {
	// TODO: make this error more general like "empty data" or something like that since it is used in cases such as AnnotationsToPromote == nil as well
	return fmt.Sprintf("InputDockerStreams is empty in CDPipeline CR named %s", err.crName)
}

func NewEmptyImageStreamErr(crName string) *EmptyImageStreamErr {
	return &EmptyImageStreamErr{crName: crName}
}

// GetInputISForStage gets InputDockerStream list by CR Stage name
func GetInputISForStage(ctx context.Context, client *k8s.RuntimeNamespacedClient, stageName, cdPipeline string) ([]string, error) {
	stageCrName := createStageCrName(cdPipeline, stageName)
	stageCR, err := client.GetCDStage(ctx, stageCrName)
	if err != nil {
		return nil, err
	}

	if stageCR.Spec.CdPipeline == "" {
		return nil, fmt.Errorf("Spec.CdPipeline is empty in Stage CR named %s", stageCrName)
	}

	cdPipelineName := stageCR.Spec.CdPipeline
	cdPipelineCR, err := client.GetCDPipeline(ctx, cdPipelineName)
	if err != nil {
		return nil, err
	}

	if stageCR.Spec.Order != 0 {
		previousStageName, errPrevStage := findPreviousStageName(stageCR.Annotations)
		if errPrevStage != nil {
			return nil, errPrevStage
		}
		return inputCISListFromApplications(ctx, client, cdPipelineCR, previousStageName)
	}

	if cdPipelineCR.Spec.InputDockerStreams == nil {
		return nil, NewEmptyImageStreamErr(cdPipelineName)
	}
	return cdPipelineCR.Spec.InputDockerStreams, nil
}

// inputCISListFromApplications gets input CodebaseImageStream list from CDPipeline and Stage CR
func inputCISListFromApplications(ctx context.Context, client *k8s.RuntimeNamespacedClient, cdPipelineCR *cdPipeApi.CDPipeline, previousStageName string) ([]string, error) {
	logger := applog.LoggerFromContext(ctx)

	var imageStream []string
	var cisNames []string
	re := strings.NewReplacer("/", "-", ".", "-")
	for _, inputDSName := range cdPipelineCR.Spec.InputDockerStreams {
		appName, err := applications.AppNameByInputIS(ctx, client, inputDSName)
		if err != nil {
			return nil, err
		}

		appName = re.Replace(appName)
		CISName := createCISName(cdPipelineCR.Name, previousStageName, appName)
		cisNames = append(cisNames, CISName)
	}

	for _, cisName := range cisNames {
		_, err := client.GetCodebaseImageStream(ctx, cisName)
		if err != nil {
			logger.Error("get codebase image stream failed", zap.Error(err), zap.String("cis_name", cisName))
			continue
		} else {
			imageStream = append(imageStream, cisName)
		}
	}

	if len(imageStream) != 0 {
		return imageStream, nil
	} else {
		//this error only occurs when we can't find any ImageStreamCR with cisName
		return nil, errors.New("InputIS verification failed")
	}
}

// GetOutputISForStage gets OutputIS list by CR CDPipeline name and Jenkins stage name
func GetOutputISForStage(ctx context.Context, client *k8s.RuntimeNamespacedClient, cdPipeCRName string, stageName string) ([]string, error) {
	cdPipeCR, err := client.GetCDPipeline(ctx, cdPipeCRName)
	if err != nil {
		return nil, err
	}
	var outputIS []string
	for _, inputDSName := range cdPipeCR.Spec.InputDockerStreams {
		appName, errApp := applications.AppNameByInputIS(ctx, client, inputDSName)
		if errApp != nil {
			return nil, errApp
		}
		outputIS = append(outputIS, createCISName(cdPipeCRName, stageName, appName))
	}
	return outputIS, nil
}

// createCISName: CIS is abbreviation for CodebaseImageStream
func createCISName(pipelineName, stageName, codebaseName string) string {
	return fmt.Sprintf("%s-%s-%s-%s", pipelineName, stageName, codebaseName, verifiedImageStreamSuffix)
}

func findPreviousStageName(annotations map[string]string) (string, error) {
	if annotations == nil {
		return "", fmt.Errorf("there is no annotation")
	}

	if val, ok := annotations[PreviousStageNameAnnotationKey]; ok {
		return val, nil
	}

	return "", fmt.Errorf("stage doesn`t contain %s annotation", PreviousStageNameAnnotationKey)
}

// createStageCrName: creates a full stageName for stage CR
func createStageCrName(cdPipelineName, stageName string) string {
	return fmt.Sprintf("%s-%s", cdPipelineName, stageName)
}

func OutputCBStreamsForCodebaseNames(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, codebaseNames []string) (map[string][]codeBaseApi.CodebaseImageStream, error) {
	cbStreams, err := k8sClient.CodebaseImageStreamList(ctx)
	if err != nil {
		return nil, err
	}

	cdNamesMap := make(map[string]struct{})
	for _, name := range codebaseNames {
		cdNamesMap[name] = struct{}{}
	}

	codebaseStreams := make(map[string][]codeBaseApi.CodebaseImageStream)
	for _, cbStream := range cbStreams {
		codebaseName := cbStream.Spec.Codebase
		if _, ok := cdNamesMap[codebaseName]; ok {
			codebaseStreams[codebaseName] = append(codebaseStreams[codebaseName], cbStream)
		}
	}
	return codebaseStreams, nil
}
