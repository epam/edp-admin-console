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

// +k8s:openapi-gen=true
type CDPipelineSpec struct {
	Name               string   `json:"name"`
	CodebaseBranch     []string `json:"codebase_branch"`
	ThirdPartyServices []string `json:"services"`
}

// +k8s:openapi-gen=true
type CDPipelineStatus struct {
	LastTimeUpdated time.Time `json:"last_time_updated"`
	Status          string    `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type CDPipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CDPipelineSpec   `json:"spec,omitempty"`
	Status CDPipelineStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CDPipelineList contains a list of CDPipeline
type CDPipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CDPipeline `json:"items"`
}
