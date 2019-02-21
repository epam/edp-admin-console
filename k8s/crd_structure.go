package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BusinessApplicationSpec struct {
	Name      string    `json:"name" valid:"Required"`
	Strategy  string    `json:"strategy" valid:"Required"`
	Lang      string    `json:"lang" valid:"Required"`
	BuildTool string    `json:"buildTool" valid:"Required"`
	Framework string    `json:"framework" valid:"Required"`
	Git       *Git      `json:"git,omitempty"`
	Route     *Route    `json:"route,omitempty"`
	Database  *Database `json:"database,omitempty"`
}

type Git struct {
	Url      string `json:"url"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
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
type BusinessApplicationStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type BusinessApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BusinessApplicationSpec   `json:"spec,omitempty"`
	Status BusinessApplicationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// impl.BusinessApplicationList contains a list of impl.BusinessApplication
type BusinessApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessApplication `json:"items"`
}
