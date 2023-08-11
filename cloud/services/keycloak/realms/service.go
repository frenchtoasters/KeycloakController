package realms

import (
	"context"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/scope"
	gokeycloak "github.com/Nerzal/gocloak/v13"
)

type realmInterface interface {
	Get(ctx context.Context, name string) (int, []*gokeycloak.RealmRepresentation, error)
}

type Scope interface {
	cloud.Keycloak
}

type Service struct {
	scope scope.KeycloakScope
	realm string
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope scope.KeycloakScope) *Service {
	return &Service{
		scope: scope,
		realm: scope.Keycloak.Spec.RealmName,
	}
}
