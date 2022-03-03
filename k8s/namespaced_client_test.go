package k8s

import (
	"context"
	"os"
	"strings"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-codebase-operator/v2/pkg/codebasebranch"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	name = "name"
	ns   = "ns"
	ns2  = "ns2"
)

func createCDBranchCRByNamespace(ns string) *codeBaseApi.CodebaseBranch {
	return &codeBaseApi.CodebaseBranch{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CodebaseBranch",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			BranchName: "master",
			Release:    true,
		},
	}
}

func createCDStageCR(ns string) *cdPipeApi.Stage {
	return &cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createCDPipelineCR(ns string) *cdPipeApi.CDPipeline {
	return &cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createCodebaseCR(ns string) *codeBaseApi.Codebase {
	return &codeBaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createCodebaseImageStreamCR(t *testing.T, ns string) *codeBaseApi.CodebaseImageStream {
	t.Helper()
	return &codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func TestRuntimeNamespacedClient_CreateCodebaseByCustomFields_AlreadyExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.Codebase{})

	initCR := createCodebaseCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	lang := "go"
	spec := codeBaseApi.CodebaseSpec{Lang: lang}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCodebaseByCustomFields(ctx, name, spec, codeBaseApi.CodebaseStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	codebase, err := k8sClient.GetCodebase(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, codebase.Name)
}

func TestRuntimeNamespacedClient_CreateCodebaseByCustomFields(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.Codebase{})

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	lang := "go"
	spec := codeBaseApi.CodebaseSpec{Lang: lang}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCodebaseByCustomFields(ctx, name, spec, codeBaseApi.CodebaseStatus{})
	assert.NoError(t, err)

	codebase, err := k8sClient.GetCodebase(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, codebase.Name)
	assert.Equal(t, lang, codebase.Spec.Lang)
}

func TestK8SClient_GetCBBranch(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cbBranchCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, "ns")
	if err != nil {
		t.Fatal(err)
	}
	branch, err := k8sClient.GetCBBranch(ctx, "name")
	assert.NoError(t, err)
	assert.Equal(t, name, branch.Name)
	assert.Equal(t, ns, branch.Namespace)
}

func TestK8SClient_GetCBBranch_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createCDBranchCRByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, branch)
}

func TestNamespacedClient_DeleteCBBranch_NotFound(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createCDBranchCRByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.DeleteCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_DeleteCBBranch(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.DeleteCBBranch(ctx, name)
	assert.NoError(t, err)

	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, branch)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.NoError(t, err)
	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, spec.Release, branch.Spec.Release)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_CreateCBBranchByCustomFields(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, "master", branch.Spec.BranchName)
}

func TestNamespacedClient_GetCBBranch_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	CR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(CR).Build()
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	branch, err := k8sClient.GetCBBranch(ctx, "name")
	assert.Nil(t, branch)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}

	err := k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_DeleteCBBranch_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cbBranchCR).Build()
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	err := k8sClient.DeleteCBBranch(ctx, "name")
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_CreateCBBranchByCustomFields_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createCDBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	err := k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestK8SClient_GetCDStage_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient := RuntimeNamespacedClient{Client: client}
	stage, err := k8sClient.GetCDStage(ctx, name)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, stage)
}

func TestK8SClient_GetCDStage_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})
	cdStageCR := createCDStageCR(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdStageCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, stage)
}

func TestK8SClient_GetCDStage(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})
	cdStageCR := createCDStageCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdStageCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, stage.Name)
	assert.Equal(t, ns, stage.Namespace)
}

func TestK8SClient_GetCDPipeline_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient := RuntimeNamespacedClient{Client: client}
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, cdPipeline)
}

func TestK8SClient_GetCDPipeline_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})
	cdPipelineCR := createCDPipelineCR(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdPipelineCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, cdPipeline)
}

func TestK8SClient_GetCDPipeline(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})
	cdPipelineCR := createCDPipelineCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdPipelineCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, cdPipeline.Name)
	assert.Equal(t, ns, cdPipeline.Namespace)
}

func TestK8SClient_GetCodebaseImageStream_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient := RuntimeNamespacedClient{Client: client}
	codebaseImageStream, err := k8sClient.GetCodebaseImageStream(ctx, name)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, codebaseImageStream)
}

func TestK8SClient_GetCodebaseImageStream_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseImageStream{})
	codebaseImageStreamCR := createCodebaseImageStreamCR(t, ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseImageStreamCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	codebaseImageStream, err := k8sClient.GetCodebaseImageStream(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, codebaseImageStream)
}

func TestK8SClient_GetCodebaseImageStream_Success(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.CodebaseImageStream{})
	codebaseImageStreamCR := createCodebaseImageStreamCR(t, ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseImageStreamCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	codebaseImageStream, err := k8sClient.GetCodebaseImageStream(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, codebaseImageStream.Name)
	assert.Equal(t, ns, codebaseImageStream.Namespace)
}

func TestK8SClient_GetCodebase_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.Codebase{})
	codebaseCR := createCodebaseCR(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	codebase, err := k8sClient.GetCodebase(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, codebase)
}

func TestSetupNamespacedClient_EnvErr(t *testing.T) {
	err := os.Unsetenv(NamespaceEnv)
	if err != nil {
		t.Fatal()
	}
	client, err := SetupNamespacedClient()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "cant find NAMESPACE env"))
	assert.Nil(t, client)
}

func TestSetupNamespacedClient_NewClientErr(t *testing.T) {
	err := os.Setenv(NamespaceEnv, "test")
	if err != nil {
		t.Fatal(err)
	}
	client, err := SetupNamespacedClient()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "cant setup client"))
	assert.Nil(t, client)
	err = os.Unsetenv(NamespaceEnv)
	if err != nil {
		t.Fatal()
	}
}

func TestK8SClient_GetCodebase(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	codebaseCR := createCodebaseCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	codebase, err := k8sClient.GetCodebase(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, codebase.Name)
	assert.Equal(t, ns, codebase.Namespace)
}

func TestNewRuntimeNamespacedClient_EmptyNamespace(t *testing.T) {
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, "")
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, k8sClient)
}

func TestRuntimeNamespacedClient_CodebaseBranchesListByCodebaseName(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	namespace_1 := "k8s"
	namespace_2 := "linux"
	crCodebaseBranch_1 := &codeBaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "docker-master",
			Namespace: namespace_1,
			Labels: map[string]string{
				codebasebranch.LabelCodebaseName: "docker",
			},
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			CodebaseName: "docker",
			BranchName:   "master",
		},
	}
	crCodebaseBranch_2 := &codeBaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "docker-master",
			Namespace: namespace_2,
			Labels: map[string]string{
				codebasebranch.LabelCodebaseName: "docker",
			},
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			CodebaseName: "docker",
			BranchName:   "master",
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(crCodebaseBranch_1, crCodebaseBranch_2).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, namespace_1)
	if err != nil {
		t.Fatal(err)
	}
	cbBranchesList, err := k8sClient.CodebaseBranchesListByCodebaseName(ctx, "docker")
	assert.NoError(t, err)
	expectedCbBranches := []*codeBaseApi.CodebaseBranch{
		crCodebaseBranch_1,
	}
	assert.Equal(t, expectedCbBranches, cbBranchesList)
}
