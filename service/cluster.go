/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"edp-admin-console/k8s"
	"edp-admin-console/service/logger"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logger.GetLogger()

type ClusterService struct {
	Clients k8s.ClientSet
}

func (this ClusterService) GetAllStorageClasses() ([]string, error) {
	var storageClasses []string
	storageClient := this.Clients.StorageClient

	classList, err := storageClient.StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error("An error has occurred while getting storage classes", zap.Error(err))
		return nil, err
	}

	if len(classList.Items) == 0 {
		return []string{}, nil
	}

	for _, element := range classList.Items {
		storageClasses = append(storageClasses, element.ObjectMeta.Name)
	}

	return storageClasses, nil
}
