package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const NamespaceEnv = "NAMESPACE"

func SetupNamespacedClient() (*RuntimeNamespacedClient, error) {
	utilRuntime.Must(codeBaseApi.AddToScheme(scheme.Scheme))
	utilRuntime.Must(cdPipeApi.AddToScheme(scheme.Scheme))
	namespace, ok := os.LookupEnv(NamespaceEnv)
	if !ok {
		return nil, errors.New("cant find NAMESPACE env")
	}
	client, err := runtimeClient.New(config.GetConfigOrDie(), runtimeClient.Options{})
	if err != nil {
		return nil, fmt.Errorf("%w cant setup client", err)
	}

	return NewRuntimeNamespacedClient(client, namespace), nil
}

type RuntimeNamespacedClient struct {
	runtimeClient.Client
	Namespace string
}

type EmptyNamespaceErr struct {
	msg string
}

func (err *EmptyNamespaceErr) Error() string {
	return err.msg
}

func NemEmptyNamespaceErr(msg string) *EmptyNamespaceErr {
	return &EmptyNamespaceErr{msg: msg}
}

func AsEmptyNamespaceErr(err error) bool {
	var emptyNamespaceErr *EmptyNamespaceErr
	return errors.As(err, &emptyNamespaceErr)
}

// NewRuntimeNamespacedClient wraps an existing client enforcing the namespace value.
func NewRuntimeNamespacedClient(client runtimeClient.Client, namespace string) *RuntimeNamespacedClient {
	return &RuntimeNamespacedClient{
		Client:    client,
		Namespace: namespace,
	}
}

// GetCBBranch retrieves an CodebaseBranch structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetCBBranch(ctx context.Context, crName string) (*codeBaseApi.CodebaseBranch, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	codebaseBranch := &codeBaseApi.CodebaseBranch{}
	nsn := types.NamespacedName{
		Name:      crName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, codebaseBranch)
	if err != nil {
		return nil, fmt.Errorf("%w. failed to get CodebaseBranch CR with crName: %s", err, crName)
	}
	return codebaseBranch, nil
}

// UpdateCBBranchByCustomFields updates only the custom fields of the CodebaseBranch CR.
func (c *RuntimeNamespacedClient) UpdateCBBranchByCustomFields(ctx context.Context, crName string, spec codeBaseApi.CodebaseBranchSpec, status codeBaseApi.CodebaseBranchStatus) error {
	codebaseBranch, err := c.GetCBBranch(ctx, crName)
	if err != nil {
		return err
	}
	codebaseBranch.Status = status
	codebaseBranch.Spec = spec
	err = c.Update(ctx, codebaseBranch)
	return err
}

// CreateCBBranchByCustomFields creates CodebaseBranch CR by custom fields and name
func (c *RuntimeNamespacedClient) CreateCBBranchByCustomFields(ctx context.Context, crName string, spec codeBaseApi.CodebaseBranchSpec, status codeBaseApi.CodebaseBranchStatus) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	codebaseBranch := &codeBaseApi.CodebaseBranch{
		ObjectMeta: v1.ObjectMeta{
			Name:      crName,
			Namespace: c.Namespace,
		},
		Spec:   spec,
		Status: status,
	}
	err := c.Create(ctx, codebaseBranch)
	return err
}

// DeleteCBBranch deletes CodebaseBranch CR from the Kubernetes Cluster by name.
func (c *RuntimeNamespacedClient) DeleteCBBranch(ctx context.Context, crName string) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	codebaseBranch := codeBaseApi.CodebaseBranch{
		ObjectMeta: v1.ObjectMeta{
			Name:      crName,
			Namespace: c.Namespace,
		},
	}
	err := c.Delete(ctx, &codebaseBranch)
	return err
}

// GetCDStage retrieves a Stage structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetCDStage(ctx context.Context, crName string) (*cdPipeApi.Stage, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	stage := &cdPipeApi.Stage{}
	nsn := types.NamespacedName{
		Name:      crName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, stage)
	if err != nil {
		return nil, err
	}
	return stage, err
}

// GetCDPipeline retrieves a CDPipeline structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetCDPipeline(ctx context.Context, crName string) (*cdPipeApi.CDPipeline, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	cdPipeline := &cdPipeApi.CDPipeline{}
	nsn := types.NamespacedName{
		Name:      crName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, cdPipeline)
	if err != nil {
		return nil, err
	}
	return cdPipeline, err
}

// GetCodebaseImageStream retrieves a CodebaseImageStream structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetCodebaseImageStream(ctx context.Context, crName string) (*codeBaseApi.CodebaseImageStream, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	codebaseImageStream := &codeBaseApi.CodebaseImageStream{}
	nsn := types.NamespacedName{
		Name:      crName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, codebaseImageStream)
	if err != nil {
		return nil, err
	}
	return codebaseImageStream, err
}

// GetCodebase retrieves a Codebase structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetCodebase(ctx context.Context, crName string) (*codeBaseApi.Codebase, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	codebase := &codeBaseApi.Codebase{}
	nsn := types.NamespacedName{
		Name:      crName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, codebase)
	if err != nil {
		return nil, err
	}
	return codebase, err
}
