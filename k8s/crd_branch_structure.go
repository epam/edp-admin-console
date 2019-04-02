package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ApplicationBranchSpec struct {
	Name            string `json:"branchName"`
	Commit          string `json:"fromCommit"`
	ApplicationName string `json:"appName"`
}

// +k8s:openapi-gen=true
type ApplicationBranchStatus struct {
	LastTimeUpdated string `json:"last_time_updated"`
	Status          string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type ApplicationBranch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationBranchSpec   `json:"spec,omitempty"`
	Status ApplicationBranchStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// impl.ApplicationBranchList contains a list of impl.ApplicationBranch
type ApplicationBranchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationBranch `json:"items"`
}
