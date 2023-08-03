package scope

import (
	"context"

	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeycloakScopeParams struct {
	Client   client.Client
	Keycloak *appdatv1alpha1.Keycloak
}

// NewKeycloakScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewKeycloakScope(ctx context.Context, params KeycloakScopeParams) (*KeycloakScope, error) {
	if params.Keycloak == nil {
		return nil, errors.New("failed to generate new scope from nil Keycloak")
	}
	if params.Keycloak == nil {
		return nil, errors.New("failed to generate new scope from nil GCPKeycloak")
	}

	helper, err := patch.NewHelper(params.Keycloak, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &KeycloakScope{
		client:      params.Client,
		Keycloak:    params.Keycloak,
		patchHelper: helper,
	}, nil
}

// KeycloakScope defines the basic context for an actuator to operate upon.
type KeycloakScope struct {
	client      client.Client
	patchHelper *patch.Helper

	Keycloak *appdatv1alpha1.Keycloak
}

// PatchObject persists the keycloak configuration and status.
func (s *KeycloakScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.Keycloak)
}

// Close closes the current scope persisting the keycloak configuration and status.
func (s *KeycloakScope) Close() error {
	return s.PatchObject()
}

func (s *KeycloakScope) Groups() []appdatv1alpha1.KeycloakGroup {
	return s.Keycloak.Spec.Groups
}

func (s *KeycloakScope) IdentityProviderRoleMapper() []appdatv1alpha1.IdentityProviderRoleMapper {
	return s.Keycloak.Spec.IdentityProviderRoleMappers
}

func (s *KeycloakScope) Users() []appdatv1alpha1.KeycloakUser {
	return s.Keycloak.Spec.Users
}

func (s *KeycloakScope) RealmName() string {
	return s.Keycloak.Spec.RealmName
}

func (s *KeycloakScope) Namespace() string {
	return s.Keycloak.Namespace
}
