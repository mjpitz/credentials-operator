package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Credentials
type Credentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CredentialsSpec   `json:"spec"`
	Status CredentialsStatus `json:"status"`
}

// Requirements
type Requirements struct {
	Length *uint `json:"length"`
	CharacterSet string `json:"characterSet"`
}

// Credential
type Credential struct {
	Key string `json:"key"`
	Requirements Requirements `json:"requirements"`
}

// SecretRef
type SecretRef struct {
	Name string `json:"name"`
}

// View
type View struct {
	SecretRef SecretRef `json:"secretRef"`
	StringDataTemplate map[string]string `json:"stringDataTemplate"`
}

// CredentialsSpec
type CredentialsSpec struct {
	Credentials []Credential `json:"credentials"`
	Views []View `json:"views"`
}

// CredentialsStatus
type CredentialsStatus struct {
}
