package users

import (
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/scope"
	gokeycloak "github.com/Nerzal/gocloak/v13"
)

type Scope interface {
	cloud.Keycloak
}

type Service struct {
	scope scope.KeycloakScope
	users []*gokeycloak.User
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope scope.KeycloakScope) *Service {
	return &Service{
		scope: scope,
		users: scope.Keycloak.Spec.Users,
	}
}
