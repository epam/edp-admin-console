package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var k8sConfig clientcmd.ClientConfig
var SchemeGroupVersion = schema.GroupVersion{Group: "edp.epam.com", Version: "v1alpha1"}

type ClientSet struct {
	CoreClient        *coreV1Client.CoreV1Client
	StorageClient     *rest.RESTClient
	ApplicationClient *rest.RESTClient
}

func init() {
	k8sConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
}

func CreateOpenShiftClients() ClientSet {
	coreClient, err := getCoreClient()
	if err != nil {
		log.Printf("An error has occurred while getting core client: %s", err)
		panic(err)
	}

	crClient, err := getApplicationClient()
	if err != nil {
		log.Printf("An error has occurred while getting custom resource client: %s", err)
		panic(err)
	}

	return ClientSet{
		CoreClient:        coreClient,
		StorageClient:     nil,
		ApplicationClient: crClient,
	}
}

func getCoreClient() (*coreV1Client.CoreV1Client, error) {
	restConfig, err := k8sConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	coreClient, err := coreV1Client.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return coreClient, nil
}

func getApplicationClient() (*rest.RESTClient, error) {
	var config *rest.Config
	var err error

	config, err = k8sConfig.ClientConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := createCrdClient(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func createCrdClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&BusinessApplication{},
		&BusinessApplicationList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
