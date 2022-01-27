package imagestream

import (
	"context"
	"errors"
	"fmt"

	"edp-admin-console/k8s"
)

type EmptyImageStreamErr struct {
	crName string
}

func (err *EmptyImageStreamErr) Error() string {
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
func GetInputISForStage(ctx context.Context, client *k8s.NamespacedClient, crName string) ([]string, error) {
	stageCR, err := client.GetCDStage(ctx, crName)
	if err != nil {
		return nil, err
	}

	if stageCR.Spec.Order != 0 {
		return nil, fmt.Errorf("this stage named %s is not the first stage", crName)
	}
	if stageCR.Spec.CdPipeline == "" {
		return nil, fmt.Errorf("Spec.CdPipeline is empty in Stage CR named %s", crName)
	}

	cdPipelineName := stageCR.Spec.CdPipeline
	cdPipelineCR, err := client.GetCDPipeline(ctx, cdPipelineName)
	if err != nil {
		return nil, err
	}

	if cdPipelineCR.Spec.InputDockerStreams == nil {
		return nil, NewEmptyImageStreamErr(cdPipelineName)
	}
	return cdPipelineCR.Spec.InputDockerStreams, nil
}
