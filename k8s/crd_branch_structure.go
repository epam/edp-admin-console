package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type CodebaseBranchSpec struct {
	Name         string `json:"branchName"`
	Commit       string `json:"fromCommit"`
	CodebaseName string `json:"codebaseName"`
}

// +k8s:openapi-gen=true
type CodebaseBranchStatus struct {
	LastTimeUpdated time.Time `json:"last_time_updated"`
	Status          string    `json:"status"`
	Username        string    `json:"username"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type CodebaseBranch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CodebaseBranchSpec   `json:"spec,omitempty"`
	Status CodebaseBranchStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// impl.CodebaseBranchList contains a list of impl.CodebaseBranch
type CodebaseBranchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CodebaseBranch `json:"items"`
}
