package users

import (
	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/scope"
)

type Scope interface {
	cloud.Keycloak
}

type Service struct {
	scope scope.KeycloakScope
	users []*appdatv1alpha1.User
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope scope.KeycloakScope) *Service {
	return &Service{
		scope: scope,
		users: scope.Keycloak.Spec.Users,
	}
}
