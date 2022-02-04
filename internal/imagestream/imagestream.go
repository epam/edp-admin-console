package imagestream

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"

	"edp-admin-console/k8s"
)

const (
	verifiedImageStreamSuffix      = "verified"
	previousStageNameAnnotationKey = "deploy.edp.epam.com/previous-stage-name"
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

func AsEmptyImageStreamErr(err error) bool {
	var emptyImageStreamEr *EmptyImageStreamErr
	return errors.As(err, &emptyImageStreamEr)
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
		previousStageName, err := findPreviousStageName(stageCR.Annotations)
		if err != nil {
			return nil, err
		}
		return inputCISListFromApplicationsToPromote(ctx, client, cdPipelineCR, previousStageName)
	}

	if cdPipelineCR.Spec.InputDockerStreams == nil {
		return nil, NewEmptyImageStreamErr(cdPipelineName)
	}
	return cdPipelineCR.Spec.InputDockerStreams, nil
}

// inputCISListFromApplicationsToPromote gets ApplicationsToPromote list by CR Stage name
func inputCISListFromApplicationsToPromote(ctx context.Context, client *k8s.RuntimeNamespacedClient, cdPipelineCR *cdPipeApi.CDPipeline, previousStageName string) ([]string, error) {
	if cdPipelineCR.Spec.ApplicationsToPromote == nil {
		return nil, NewEmptyImageStreamErr(cdPipelineCR.Name)
	}

	var imageStream []string
	var cisNames []string

	for _, name := range cdPipelineCR.Spec.ApplicationsToPromote {
		re := strings.NewReplacer("/", "-", ".", "-")
		name = re.Replace(name)
		CISName := createCISName(cdPipelineCR.Name, previousStageName, name)
		cisNames = append(cisNames, CISName)
	}

	for _, cisName := range cisNames {
		_, err := client.GetCodebaseImageStream(ctx, cisName)
		if err != nil {
			//TODO: add error logging
			continue
		} else {
			imageStream = append(imageStream, cisName)
		}
	}

	if len(imageStream) != 0 {
		return imageStream, nil
	} else {
		return nil, NewEmptyImageStreamErr(cdPipelineCR.Name)
	}
}

// createCISName: CIS is abbreviation for CodebaseImageStream
func createCISName(pipelineName, stageName, codebaseName string) string {
	return fmt.Sprintf("%s-%s-%s-%s", pipelineName, stageName, codebaseName, verifiedImageStreamSuffix)
}

func findPreviousStageName(annotations map[string]string) (string, error) {
	if annotations == nil {
		return "", fmt.Errorf("there is no annotation")
	}

	if val, ok := annotations[previousStageNameAnnotationKey]; ok {
		return val, nil
	}

	return "", fmt.Errorf("stage doesn`t contain %s annotation", previousStageNameAnnotationKey)
}

// createStageCrName: creates a full stageName for stage CR
func createStageCrName(cdPipelineName, stageName string) string {
	return fmt.Sprintf("%s-%s", cdPipelineName, stageName)
}
