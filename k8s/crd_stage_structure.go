package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// +k8s:openapi-gen=true
type StageSpec struct {
	Name        string `json:"name"`
	CdPipeline  string `json:"cdPipeline"`
	Description string `json:"description"`
	QualityGate string `json:"qualityGate"`
	TriggerType string `json:"triggerType"`
	Input       string `json:"input"`
}

// +k8s:openapi-gen=true
type StageStatus struct {
	LastTimeUpdated time.Time `json:"last_time_updated"`
	Status          string    `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type Stage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StageSpec   `json:"spec,omitempty"`
	Status StageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StageList contains a list of Stage
type StageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stage `json:"items"`
}
