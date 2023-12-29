package scope

import (
	"context"

	gokeycloak "github.com/Nerzal/gocloak/v13"

	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeycloakScopeParams struct {
	Client              client.Client
	Keycloak            *appdatv1alpha1.Keycloak
	KeycloakInstanceUrl string
	KeycloakAdminUser   string
	KeycloakAdminPass   string
	KeycloakRealmName   string
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

	client := gokeycloak.NewClient(params.KeycloakInstanceUrl)
	token, err := client.LoginAdmin(ctx, params.KeycloakAdminUser, params.KeycloakAdminPass, "master")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create admin token")
	}

	return &KeycloakScope{
		client:         params.Client,
		patchHelper:    helper,
		KeycloakToken:  token,
		KeycloakClient: client,
		Keycloak:       params.Keycloak,
	}, nil
}

// KeycloakScope defines the basic context for an actuator to operate upon.
type KeycloakScope struct {
	client         client.Client
	patchHelper    *patch.Helper
	KeycloakToken  *gokeycloak.JWT
	KeycloakClient *gokeycloak.GoCloak

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

func (s *KeycloakScope) Groups() []*appdatv1alpha1.Group {
	return s.Keycloak.Spec.Groups
}

func (s *KeycloakScope) IdentityProviderRoleMapper() []*appdatv1alpha1.IdentityProviderMapper {
	return s.Keycloak.Spec.IdentityProviderRoleMappers
}

func (s *KeycloakScope) Users() []*appdatv1alpha1.User {
	return s.Keycloak.Spec.Users
}

func (s *KeycloakScope) User(i int) *appdatv1alpha1.User {
	return s.Keycloak.Spec.Users[i]
}

func (s *KeycloakScope) RealmName() string {
	return s.Keycloak.Spec.RealmName
}

func (s *KeycloakScope) Namespace() string {
	return s.Keycloak.ObjectMeta.Namespace
}

func (s *KeycloakScope) Realms() string {
	return s.Keycloak.Spec.RealmName
}

func (s *KeycloakScope) Token() string {
	return s.KeycloakToken.AccessToken
}
