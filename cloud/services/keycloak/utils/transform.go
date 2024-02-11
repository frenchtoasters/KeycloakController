package utils

import (
	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/Nerzal/gocloak/v13"
	"github.com/jinzhu/copier"
)

func UserTransform(user *appdatv1alpha1.User) gocloak.User {
	gocloakUser := &gocloak.User{}
	copier.Copy(&gocloakUser, user)
	return *gocloakUser
}

func GroupTransform(group *appdatv1alpha1.Group) gocloak.Group {
	gocloakGroup := &gocloak.Group{}
	copier.Copy(&gocloakGroup, group)
	return *gocloakGroup
}

func RoleTransform(role *gocloak.Role) *appdatv1alpha1.Role {
	appdatrole := &appdatv1alpha1.Role{}
	copier.Copy(&appdatrole, role)
	return appdatrole
}
