package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Credential
type Credential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CredentialSpec   `json:"spec"`
	Status CredentialStatus `json:"status"`
}

// Requirements
type Requirements struct {
	Length       uint   `json:"length"`
	CharacterSet string `json:"characterSet"`
}

// CredentialRequirements
type CredentialRequirements struct {
	Key          string       `json:"key"`
	Requirements Requirements `json:"requirements"`
}

// SecretRef
type SecretRef struct {
	Name string `json:"name"`
}

// CredentialView
type CredentialView struct {
	SecretRef          SecretRef         `json:"secretRef"`
	StringDataTemplate map[string]string `json:"stringDataTemplate"`
}

// CredentialSpec
type CredentialSpec struct {
	Credentials []CredentialRequirements `json:"credentials"`
	Views       []CredentialView         `json:"views"`
}

// CredentialStatus
type CredentialStatus struct {
}
