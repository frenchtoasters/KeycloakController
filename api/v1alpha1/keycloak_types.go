/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PausedAnnotation is an annotation that can be applied to any MriKeycloak
	// object to prevent a controller from processing a resource.
	//
	// Controllers working with MriKeycloak API objects must check the existence of this annotation
	// on the reconciled object.
	PausedAnnotation = "mrikeycloak.x-k8s.io/paused"

	// KeycloakFinalizer is the finalizer used by the cluster controller to
	// cleanup the cluster resources when a Keycloak is being deleted.
	KeycloakFinalizer = "keycloak.mrikeycloak.x-k8s.io"
)

type KeycloakUser struct {
	FristName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Email     string          `json:"email"`
	Uupic     string          `json:"uupic"`
	Groups    []KeycloakGroup `json:"groups"`
	Roles     []KeycloakRole  `json:"roles"`
}

type KeycloakGroup struct {
	Name string `json:"name"`
}

type KeycloakRole struct {
	Name string `json:"name"`
}

type IdentityProviderRoleMapper struct {
	Provider string `json:"provider"`
	Role     string `json:"role"`
}

// KeycloakSpec defines the desired state of Keycloak
type KeycloakSpec struct {
	RealmName                   string                       `json:"realmName"`
	Users                       []KeycloakUser               `json:"users"`
	Groups                      []KeycloakGroup              `json:"groups"`
	Roles                       []KeycloakRole               `json:"roles"`
	IdentityProviderRoleMappers []IdentityProviderRoleMapper `json:"identityProviderRoleMappers"`
	Paused                      bool                         `json:"paused"`
}

// KeycloakStatus defines the observed state of Keycloak
type KeycloakStatus struct {
	Created                     bool                         `json:"created"`
	RealmName                   string                       `json:"realmName"`
	IdentityProviderRoleMappers []IdentityProviderRoleMapper `json:"identityProviderRoleMappers"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Keycloak is the Schema for the keycloaks API
type Keycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakSpec   `json:"spec,omitempty"`
	Status KeycloakStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeycloakList contains a list of Keycloak
type KeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keycloak `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Keycloak{}, &KeycloakList{})
}
