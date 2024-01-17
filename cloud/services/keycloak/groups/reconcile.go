package groups

import (
	"context"
	"fmt"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/utils"
	gokeycloak "github.com/Nerzal/gocloak/v13"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const ()

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling group resources")

	for i := range s.groups {
		params := &gokeycloak.GetGroupsParams{
			Exact:  gokeycloak.BoolP(true),
			Search: s.groups[i].Name,
		}
		group, err := s.scope.KeycloakClient.GetGroups(ctx, s.scope.Token(), s.scope.RealmName(), *params)
		if err != nil {
			return fmt.Errorf("error getting groups: %s", err)
		}
		if len(group) != 1 {
			log.Info("Reconciling group")
			groupId, err := s.scope.KeycloakClient.CreateGroup(ctx, s.scope.Token(), s.scope.RealmName(), utils.GroupTransform(s.groups[i]))
			if err != nil {
				return fmt.Errorf("error unable to create group[%s]: %s", *s.groups[i].Name, err)
			}
			log.Info(fmt.Sprintf("Created Group - [%s]", groupId))
			s.groups[i].ID = &groupId
			continue
		}
	}

	realmGroups, err := s.scope.KeycloakClient.GetGroups(ctx, s.scope.Token(), s.scope.RealmName(), gokeycloak.GetGroupsParams{})
	if err != nil {
		return fmt.Errorf("error getting groups: %s", err)
	}

	if len(s.groups) != len(realmGroups) {
		log.Info("Reconciling missing groups")
		specGroupIds := utils.ParseSpecGroupIds(s.groups)
		for i := range realmGroups {
			if slices.Contains(specGroupIds, *realmGroups[i].ID) {
				continue
			}
			err := s.scope.KeycloakClient.DeleteGroup(ctx, s.scope.Token(), s.scope.RealmName(), *realmGroups[i].ID)
			if err != nil {
				log.Info(fmt.Sprintf("unable to delete group - %s", err))
			}
			log.Info(fmt.Sprintf("Deleted group - [%s]", *realmGroups[i].ID))
		}
	}

	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting group resources")
	for i := range s.groups {
		err := s.scope.KeycloakClient.DeleteGroup(ctx, s.scope.Token(), s.scope.RealmName(), *s.groups[i].ID)
		if err != nil {
			log.Info(fmt.Sprintf("Unable to delete user[%s] - %s", *s.groups[i].ID, err))
		}
		log.Info(fmt.Sprintf("Deleted group - [%s]", *s.scope.Keycloak.Spec.Groups[i].Name))
	}

	return nil
}
