package applications

import (
	"context"
	"fmt"
	"testing"

	"edp-admin-console/k8s"
	"edp-admin-console/util/consts"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-codebase-operator/v2/pkg/codebasebranch"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func createCodebaseCRWithOptions(opts ...CodebaseCROption) *codeBaseApi.Codebase {
	codebaseCR := new(codeBaseApi.Codebase)
	for i := range opts {
		opts[i](codebaseCR)
	}
	return codebaseCR
}

type CodebaseCROption func(codebase *codeBaseApi.Codebase)

func WithCBCrName(crName string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithCBCrNamespace(crNamespace string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Namespace = crNamespace
	}
}

func WithCBCrStatus(status string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Status.Value = status
	}
}

func WithCBCrType(crType string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Type = crType
	}
}

func WithCBCrLang(lang string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Lang = lang
	}
}

func TestCodebasesByTypeAndStatus_OK(t *testing.T) {
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	crStatus_1 := consts.ActiveValue
	crType_1 := consts.Application
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
		WithCBCrStatus(crStatus_1),
		WithCBCrType(crType_1),
	)

	expectedCBCR := []codeBaseApi.Codebase{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName_1,
				Namespace: namespace,
			},
			Status: codeBaseApi.CodebaseStatus{
				Value: consts.ActiveValue,
			},
			Spec: codeBaseApi.CodebaseSpec{
				Type: consts.Application,
			},
		},
	}
	cbList := []codeBaseApi.Codebase{*codebaseCR_1}
	crStatus := consts.ActiveValue
	crType := consts.Application
	filteredList, err := CodebasesByTypeAndStatus(cbList, crType, crStatus)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, filteredList)
}

type CodebaseBranchCROption func(codebase *codeBaseApi.CodebaseBranch)

func createCodebaseBranchCRWithOptions(opts ...CodebaseBranchCROption) *codeBaseApi.CodebaseBranch {
	codebaseBranchCR := new(codeBaseApi.CodebaseBranch)
	for i := range opts {
		opts[i](codebaseBranchCR)
	}
	return codebaseBranchCR
}

func cbBranchWithName(name string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Name = name
	}
}

func cbBranchWithNamespace(namespace string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Namespace = namespace
	}
}

func cbBranchWithLabels(labels map[string]string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Labels = labels
	}
}

func cbBranchWithSpecBranchName(branchName string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Spec.BranchName = branchName
	}
}

func cbBranchWithStatus(status string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Status.Value = status
	}
}

func TestActiveCodebaseBranches_OK(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	cbCrName_1 := "cb_1"
	crStatus_1 := consts.ActiveValue
	crType_1 := consts.Application
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(cbCrName_1),
		WithCBCrNamespace(namespace),
		WithCBCrStatus(crStatus_1),
		WithCBCrType(crType_1),
	)

	cbBranchName_1 := "develop"
	crCbBranchName := fmt.Sprintf("%s-%s", cbCrName_1, cbBranchName_1)
	cbLabels_1 := map[string]string{
		codebasebranch.LabelCodebaseName: cbCrName_1,
	}
	cbCrStatus_1 := consts.ActiveValue
	stubCodebaseBranch_1 := createCodebaseBranchCRWithOptions(
		cbBranchWithName(crCbBranchName),
		cbBranchWithNamespace(namespace),
		cbBranchWithLabels(cbLabels_1),
		cbBranchWithSpecBranchName(cbBranchName_1),
		cbBranchWithStatus(cbCrStatus_1),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion,
		&codeBaseApi.Codebase{},
		&codeBaseApi.CodebaseBranch{}, &codeBaseApi.CodebaseBranchList{},
	)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR_1, stubCodebaseBranch_1).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedCBBranchesCR := []*codeBaseApi.CodebaseBranch{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:            crCbBranchName,
				Namespace:       namespace,
				Labels:          cbLabels_1,
				ResourceVersion: "999",
			},
			Spec: codeBaseApi.CodebaseBranchSpec{
				BranchName: cbBranchName_1,
			},
			Status: codeBaseApi.CodebaseBranchStatus{
				Value: cbCrStatus_1,
			},
		},
	}

	filteredList, err := ActiveCodebaseBranches(ctx, k8sClient, cbCrName_1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBBranchesCR, filteredList)
}

func TestActiveGroovyLibs_OK(t *testing.T) {
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	crStatus_1 := consts.ActiveValue
	crType_1 := consts.Library
	crLang_1 := "groovy-pipeline"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
		WithCBCrStatus(crStatus_1),
		WithCBCrType(crType_1),
		WithCBCrLang(crLang_1),
	)

	expectedCBCR := []codeBaseApi.Codebase{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName_1,
				Namespace: namespace,
			},
			Status: codeBaseApi.CodebaseStatus{
				Value: consts.ActiveValue,
			},
			Spec: codeBaseApi.CodebaseSpec{
				Type: consts.Library,
				Lang: crLang_1,
			},
		},
	}
	cbList := []codeBaseApi.Codebase{*codebaseCR_1}
	filteredList, err := ActiveGroovyLibs(cbList)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, filteredList)
}

func TestActiveAutotests_OK(t *testing.T) {
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	crStatus_1 := consts.ActiveValue
	crType_1 := consts.Autotest
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
		WithCBCrStatus(crStatus_1),
		WithCBCrType(crType_1),
	)

	expectedCBCR := []codeBaseApi.Codebase{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName_1,
				Namespace: namespace,
			},
			Status: codeBaseApi.CodebaseStatus{
				Value: consts.ActiveValue,
			},
			Spec: codeBaseApi.CodebaseSpec{
				Type: consts.Autotest,
			},
		},
	}
	cbList := []codeBaseApi.Codebase{*codebaseCR_1}
	filteredList, err := ActiveAutotests(cbList)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, filteredList)
}

func TestActiveApplications_OK(t *testing.T) {
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	crStatus_1 := consts.ActiveValue
	crType_1 := consts.Application
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
		WithCBCrStatus(crStatus_1),
		WithCBCrType(crType_1),
	)

	expectedCBCR := []codeBaseApi.Codebase{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName_1,
				Namespace: namespace,
			},
			Status: codeBaseApi.CodebaseStatus{
				Value: consts.ActiveValue,
			},
			Spec: codeBaseApi.CodebaseSpec{
				Type: consts.Application,
			},
		},
	}
	cbList := []codeBaseApi.Codebase{*codebaseCR_1}
	filteredList, err := ActiveApplications(cbList)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, filteredList)
}

func TestByNameIFExists_OK(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.Codebase{})
	fakeK8SClient := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR_1).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(fakeK8SClient, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedCBCR := &codeBaseApi.Codebase{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Codebase",
			APIVersion: "v2.edp.epam.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            crName_1,
			Namespace:       namespace,
			ResourceVersion: "999",
		},
	}
	k8sCodebase, err := ByNameIFExists(ctx, k8sClient, crName_1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, k8sCodebase)
}

func TestByNameIFExists_NotFound(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName_1 := "cb_1"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(crName_1),
		WithCBCrNamespace(namespace),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.Codebase{})
	fakeK8SClient := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR_1).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(fakeK8SClient, namespace)
	if err != nil {
		t.Fatal(err)
	}

	notFound := "not-found"
	k8sCodebase, err := ByNameIFExists(ctx, k8sClient, notFound)

	var expectedCBCR *codeBaseApi.Codebase
	assert.NoError(t, err)
	assert.Equal(t, expectedCBCR, k8sCodebase)
}
