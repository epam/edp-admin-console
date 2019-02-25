package service

import (
	"edp-admin-console/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type ClusterService struct {
	Clients k8s.ClientSet
}

func (this ClusterService) GetAllStorageClasses() ([]string, error) {
	var storageClasses []string
	storageClient := this.Clients.StorageClient

	classList, err := storageClient.StorageClasses().List(metav1.ListOptions{})
	if err != nil {
		log.Printf("An error has occurred while getting storage classes: %s", err)
		return storageClasses, err
	}

	for _, element := range classList.Items {
		storageClasses = append(storageClasses, element.ObjectMeta.Name)
	}

	return storageClasses, nil
}
