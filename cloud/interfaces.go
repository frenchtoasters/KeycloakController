package cloud

import (
	"context"

	gokeycloak "github.com/Nerzal/gocloak/v13"

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

type Keycloak interface {
	KeycloakGetter
}

// KeycloakGetter is an interface which can get cluster information.
type KeycloakGetter interface {
	Client
	RealmName() string
	Namespace() string
	Users() []*gokeycloak.User
	Groups() []*gokeycloak.Group
	IdentityProviderRoleMapper() []gokeycloak.IdentityProviderMapper
}
