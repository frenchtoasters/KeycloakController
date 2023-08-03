package cloud

import (
	"context"

	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type KeycloakApi = goKeycloakInterface

// Reconciler is a generic interface used by components offering a type of service.
type Reconciler interface {
	Reconcile(ctx context.Context) error
	Delete(ctx context.Context) error
}

// ReconcilerWithResult is a generic interface used by components offering a type of service.
type ReconcilerWithResult interface {
	Reconcile(ctx context.Context) (ctrl.Result, error)
	Delete(ctx context.Context) (ctrl.Result, error)
}

// Client is an interface which can get cloud client.
type Client interface {
	Keycloak() KeycloakApi
}

// KeycloakGetter is an interface which can get cluster information.
type KeycloakGetter interface {
	Client
	RealmName() string
	Namespace() string
	Users() []appdatv1alpha1.KeycloakUser
	Groups() []appdatv1alpha1.KeycloakGroup
	IdentityProviderRoleMapper() []appdatv1alpha1.IdentityProviderRoleMapper
}

type Keycloak interface {
	KeycloakGetter
}
