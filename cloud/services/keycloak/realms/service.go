package realms

import (
	"context"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud"
	gokeycloak "github.com/Nerzal/gocloak/v13"
)

type realmInterface interface {
	Get(ctx context.Context, name string) (int, []*gokeycloak.RealmRepresentation, error)
}

type Scope interface {
	cloud.KeycloakGetter
}

type Service struct {
	scope  Scope
	realms realmInterface
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	return &Service{
		scope:  scope,
		realms: scope.Keycloak().Realms(),
	}
}
