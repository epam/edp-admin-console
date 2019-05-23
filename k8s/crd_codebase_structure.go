/*
 * Copyright 2019 EPAM Systems.
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

package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type CodebaseSpec struct {
	Name                string      `json:"name" valid:"Required"`
	Strategy            string      `json:"strategy" valid:"Required"`
	Lang                string      `json:"lang" valid:"Required"`
	BuildTool           string      `json:"buildTool" valid:"Required"`
	Framework           string      `json:"framework" valid:"Required"`
	Repository          *Repository `json:"repository,omitempty"`
	Route               *Route      `json:"route,omitempty"`
	Database            *Database   `json:"database,omitempty"`
	TestReportFramework *string     `json:"testReportFramework"`
	Type                string      `json:"type"`
	Description         string      `json:"description"`
}

type Repository struct {
	Url string `json:"url"`
}

type Route struct {
	Site string `json:"site"`
	Path string `json:"path"`
}

type Database struct {
	Kind     string `json:"kind"`
	Version  string `json:"version"`
	Capacity string `json:"capacity"`
	Storage  string `json:"storage"`
}

// +k8s:openapi-gen=true
type CodebaseStatus struct {
	Available       bool      `json:"available"`
	LastTimeUpdated time.Time `json:"last_time_updated"`
	Status          string    `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type Codebase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CodebaseSpec   `json:"spec,omitempty"`
	Status CodebaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// impl.CodebaseList contains a list of impl.Codebase
type CodebaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Codebase `json:"items"`
}
