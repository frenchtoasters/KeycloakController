package utils

import (
	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/Nerzal/gocloak/v13"
)

func UserTransform(user *appdatv1alpha1.User) gocloak.User {
	gocloakUser := &gocloak.User{}
	gocloakUser.Access = user.Access
	gocloakUser.ID = user.ID
	gocloakUser.Totp = user.Totp
	gocloakUser.Email = user.Email
	gocloakUser.Groups = user.Groups
	gocloakUser.Enabled = user.Enabled
	gocloakUser.LastName = user.LastName
	gocloakUser.Username = user.Username
	gocloakUser.FirstName = user.FirstName
	gocloakUser.Attributes = user.Attributes
	gocloakUser.RealmRoles = user.RealmRoles
	gocloakUser.ClientRoles = user.ClientRoles
	gocloakUser.EmailVerified = user.EmailVerified
	gocloakUser.FederationLink = user.FederationLink
	gocloakUser.RequiredActions = user.RequiredActions
	gocloakUser.CreatedTimestamp = user.CreatedTimestamp
	gocloakUser.ServiceAccountClientID = user.ServiceAccountClientID

	for i := 0; i < len(*user.Credentials); i++ {
		(*gocloakUser.Credentials)[i].CreatedDate = (*user.Credentials)[i].CreatedDate
		(*gocloakUser.Credentials)[i].Temporary = (*user.Credentials)[i].Temporary
		(*gocloakUser.Credentials)[i].Type = (*user.Credentials)[i].Type
		(*gocloakUser.Credentials)[i].Value = (*user.Credentials)[i].Value
		(*gocloakUser.Credentials)[i].CredentialData = (*user.Credentials)[i].CredentialData
		(*gocloakUser.Credentials)[i].ID = (*user.Credentials)[i].ID
		(*gocloakUser.Credentials)[i].Priority = (*user.Credentials)[i].Priority
		(*gocloakUser.Credentials)[i].SecretData = (*user.Credentials)[i].SecretData
		(*gocloakUser.Credentials)[i].UserLabel = (*user.Credentials)[i].UserLabel
	}
	return *gocloakUser
}
