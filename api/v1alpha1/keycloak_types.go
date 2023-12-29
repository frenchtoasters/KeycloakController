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

type User struct {
	ID               *string              `json:"id,omitempty"`
	CreatedTimestamp *int64               `json:"createdTimestamp,omitempty"`
	Username         *string              `json:"username,omitempty"`
	Enabled          *bool                `json:"enabled,omitempty"`
	Totp             *bool                `json:"totp,omitempty"`
	EmailVerified    *bool                `json:"emailVerified,omitempty"`
	FirstName        *string              `json:"firstName,omitempty"`
	LastName         *string              `json:"lastName,omitempty"`
	Email            *string              `json:"email,omitempty"`
	FederationLink   *string              `json:"federationLink,omitempty"`
	Attributes       *map[string][]string `json:"attributes,omitempty"`
	// Currently cannot be used with controller-gen because of the []interface{}
	// DisableableCredentialTypes *[]interface{}              `json:"disableableCredentialTypes,omitempty"`
	RequiredActions        *[]string                   `json:"requiredActions,omitempty"`
	Access                 *map[string]bool            `json:"access,omitempty"`
	ClientRoles            *map[string][]string        `json:"clientRoles,omitempty"`
	RealmRoles             *[]string                   `json:"realmRoles,omitempty"`
	Groups                 *[]string                   `json:"groups,omitempty"`
	ServiceAccountClientID *string                     `json:"serviceAccountClientId,omitempty"`
	Credentials            *[]CredentialRepresentation `json:"credentials,omitempty"`
}

// CredentialRepresentation is a representations of the credentials
// https://www.keycloak.org/docs-api/22.0.1/rest-api/index.html#CredentialRepresentation
type CredentialRepresentation struct {
	CreatedDate    *int64  `json:"createdDate,omitempty"`
	Temporary      *bool   `json:"temporary,omitempty"`
	Type           *string `json:"type,omitempty"`
	Value          *string `json:"value,omitempty"`
	CredentialData *string `json:"credentialData,omitempty"`
	ID             *string `json:"id,omitempty"`
	Priority       *int32  `json:"priority,omitempty"`
	SecretData     *string `json:"secretData,omitempty"`
	UserLabel      *string `json:"userLabel,omitempty"`
}

// Group is a Group
type Group struct {
	ID          *string              `json:"id,omitempty"`
	Type        *string              `json:"type,omitempty"`
	Name        *string              `json:"name,omitempty"`
	Path        *string              `json:"path,omitempty"`
	Attributes  *map[string][]string `json:"attributes,omitempty"`
	Access      *map[string]bool     `json:"access,omitempty"`
	ClientRoles *map[string][]string `json:"clientRoles,omitempty"`
	RealmRoles  *[]string            `json:"realmRoles,omitempty"`
}

// CompositesRepresentation represents the composite roles of a role
type CompositesRepresentation struct {
	Client *map[string][]string `json:"client,omitempty"`
	Realm  *[]string            `json:"realm,omitempty"`
}

// Role is a role
type Role struct {
	ID                 *string                   `json:"id,omitempty"`
	Name               *string                   `json:"name,omitempty"`
	ScopeParamRequired *bool                     `json:"scopeParamRequired,omitempty"`
	Composite          *bool                     `json:"composite,omitempty"`
	Composites         *CompositesRepresentation `json:"composites,omitempty"`
	ClientRole         *bool                     `json:"clientRole,omitempty"`
	ContainerID        *string                   `json:"containerId,omitempty"`
	Description        *string                   `json:"description,omitempty"`
	Attributes         *map[string][]string      `json:"attributes,omitempty"`
}

// IdentityProviderMapper represents the body of a call to add a mapper to
// an identity provider
type IdentityProviderMapper struct {
	ID                     *string            `json:"id,omitempty"`
	Name                   *string            `json:"name,omitempty"`
	IdentityProviderMapper *string            `json:"identityProviderMapper,omitempty"`
	IdentityProviderAlias  *string            `json:"identityProviderAlias,omitempty"`
	Config                 *map[string]string `json:"config"`
}

// KeycloakSpec defines the desired state of Keycloak
type KeycloakSpec struct {
	RealmName                   string                    `json:"realmName"`
	ManagedRealm                bool                      `json:"managedRealm"`
	Users                       []*User                   `json:"users"`
	Groups                      []*Group                  `json:"groups"`
	Roles                       []*Role                   `json:"roles"`
	IdentityProviderRoleMappers []*IdentityProviderMapper `json:"identityProviderRoleMappers"`
	Paused                      bool                      `json:"paused"`
}

// KeycloakStatus defines the observed state of Keycloak
type KeycloakStatus struct {
	RealmName string `json:"realmName"`
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
