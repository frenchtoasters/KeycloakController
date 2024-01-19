package utils

import (
	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/Nerzal/gocloak/v13"
)

func ListEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ParseSpecGroupIds(groups []*appdatv1alpha1.Group) []string {
	ids := []string{}
	for i := range groups {
		ids = append(ids, *groups[i].ID)
	}
	return ids
}

func ParseSpecNames(users []*appdatv1alpha1.User) []*string {
	ids := []*string{}
	for i := range users {
		ids = append(ids, users[i].Username)
	}
	return ids
}

func ParseUserGroups(groups []*gocloak.Group) []string {
	names := []string{}
	for i := range groups {
		names = append(names, *groups[i].Name)
	}
	return names
}

func Contains(userRoles []string, role *string) bool {
	for _, v := range userRoles {
		if v == *role {
			return true
		}
	}
	return false
}
