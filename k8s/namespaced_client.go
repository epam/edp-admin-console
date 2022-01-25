package k8s

import (
	"context"
	"errors"
	"fmt"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type NamespacedClient struct {
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

// NewNamespacedClient wraps an existing client enforcing the namespace value.
func NewNamespacedClient(client runtimeClient.Client, namespace string) *NamespacedClient {
	return &NamespacedClient{
		Client:    client,
		Namespace: namespace,
	}
}

// GetCBBranch retrieves an CodebaseBranch structure for the given custom resource name from the Kubernetes Cluster CR.
func (c *NamespacedClient) GetCBBranch(ctx context.Context, nameCR string) (codeBaseApi.CodebaseBranch, error) {
	if c.Namespace == "" {
		return codeBaseApi.CodebaseBranch{}, NemEmptyNamespaceErr("client namespace is not set")
	}
	instanceCB := codeBaseApi.CodebaseBranch{}
	nsn := types.NamespacedName{
		Name:      nameCR,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, nsn, &instanceCB)
	if err != nil {
		return codeBaseApi.CodebaseBranch{}, fmt.Errorf("%w. failed to get CodebaseBranch CR with nameCR: %s", err, nameCR)
	}
	return instanceCB, nil
}

// UpdateCBBranchByCustomFields updates only the custom fields of the CodebaseBranch CR.
func (c *NamespacedClient) UpdateCBBranchByCustomFields(ctx context.Context, nameCR string, spec codeBaseApi.CodebaseBranchSpec, status codeBaseApi.CodebaseBranchStatus) error {
	instance, err := c.GetCBBranch(ctx, nameCR)
	if err != nil {
		return err
	}
	instance.Status = status
	instance.Spec = spec
	err = c.Update(ctx, &instance)
	return err
}

// CreateCBBranchByCustomFields creates CodebaseBranch CR by custom fields and name
func (c *NamespacedClient) CreateCBBranchByCustomFields(ctx context.Context, nameCR string, spec codeBaseApi.CodebaseBranchSpec, status codeBaseApi.CodebaseBranchStatus) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	instance := &codeBaseApi.CodebaseBranch{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameCR,
			Namespace: c.Namespace,
		},
		Spec:   spec,
		Status: status,
	}
	err := c.Create(ctx, instance)
	return err
}

// DeleteCBBranch deletes CodebaseBranch CR from the Kubernetes Cluster by name.
func (c *NamespacedClient) DeleteCBBranch(ctx context.Context, nameCR string) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	instanceCB := codeBaseApi.CodebaseBranch{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameCR,
			Namespace: c.Namespace,
		},
	}
	err := c.Delete(ctx, &instanceCB)
	return err
}
