package users

import (
	"context"
	"fmt"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/utils"
	gokeycloak "github.com/Nerzal/gocloak/v13"
	"golang.org/x/exp/slices"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling user resources")

	for i := range s.users {
		params := &gokeycloak.GetUsersParams{
			Username: s.users[i].Username,
			Exact:    gokeycloak.BoolP(true),
		}
		user, err := s.scope.KeycloakClient.GetUsers(ctx, s.scope.Token(), s.scope.RealmName(), *params)
		if err != nil {
			return fmt.Errorf("error getting users: %s", err)
		}
		if len(user) == 0 {
			log.Info("Reconciling new users")
			userId, err := s.scope.KeycloakClient.CreateUser(ctx, s.scope.Token(), s.scope.RealmName(), utils.UserTransform(s.users[i]))
			if err != nil {
				return fmt.Errorf("error unable to create user: %s", err)
			}
			log.Info(fmt.Sprintf("Created user - [%s]", userId))
			s.users[i].ID = &userId
			continue
		}

		userGroups, err := s.scope.KeycloakClient.GetUserGroups(ctx, s.scope.Token(), s.scope.RealmName(), *user[0].ID, gokeycloak.GetGroupsParams{})
		if err != nil {
			return fmt.Errorf("error getting users groups[%s]: %s", *user[0].Username, err)
		}

		groups := utils.ParseUserGroups(userGroups)
		if !utils.ListEqual(groups, *s.users[i].Groups) {
			log.Info(fmt.Sprintf("Updating user - [%s]", *user[0].Username))
			err := s.scope.KeycloakClient.UpdateUser(ctx, s.scope.Token(), s.scope.RealmName(), utils.UserTransform(s.users[i]))
			if err != nil {
				return fmt.Errorf("error updating users groups[%s]: %s", *s.users[i].Username, err)
			}
		}
	}

	userCount, err := s.scope.KeycloakClient.GetUserCount(ctx, s.scope.Token(), s.scope.RealmName(), gokeycloak.GetUsersParams{})
	if err != nil {
		return fmt.Errorf("error getting user count: %s", err)
	}

	if userCount != len(s.users) {
		log.Info("Reconciling missing users")
		realmUsers, err := s.scope.KeycloakClient.GetUsers(ctx, s.scope.Token(), s.scope.RealmName(), gokeycloak.GetUsersParams{})
		if err != nil {
			return fmt.Errorf("error getting users: %s", err)
		}

		specUserIds := utils.ParseSpecIds(s.users)
		for i := range realmUsers {
			if slices.Contains(specUserIds, realmUsers[i].ID) {
				continue
			}
			err := s.scope.KeycloakClient.DeleteUser(ctx, s.scope.Token(), s.scope.RealmName(), *realmUsers[i].ID)
			if err != nil {
				log.Info(fmt.Sprintf("unable to delete user - %s", err))
			}
			log.Info(fmt.Sprintf("Deleted user - [%s]", *realmUsers[i].ID))
		}
	}

	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting user resources")
	// TODO:: Figure out how to handle deleting specific users from a list
	for i := range s.users {
		err := s.scope.KeycloakClient.DeleteUser(ctx, s.scope.Token(), s.scope.RealmName(), *s.users[i].ID)
		if err != nil {
			log.Info("Unable to delete user", "error", err)
		}
		log.Info("Deleted user", "user", s.users[i].Username)
	}

	return nil
}
