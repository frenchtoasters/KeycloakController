package users

import (
	"context"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/utils"
	gokeycloak "github.com/Nerzal/gocloak/v13"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling user resources")
	for i := range s.users {
		params := &gokeycloak.GetUsersParams{
			Username: s.users[i].Username,
		}
		user, err := s.scope.KeycloakClient.GetUsers(ctx, s.scope.Token(), s.scope.RealmName(), *params)
		if err != nil {
			// TODO: Check that error is recoverable here
			log.Info("Unable to find user - %s", s.users[i].Username)
			log.Info("Creating user - %s", s.users[i].Username)
			userId, err := s.scope.KeycloakClient.CreateUser(ctx, s.scope.Token(), s.scope.RealmName(), utils.UserTransform(s.users[i]))
			if err != nil {
				log.Info("Unable to create user - %v", err)
			}
			log.Info("Created user[%s] - %s", s.users[i].Username, userId)
			continue
		}
		log.Info("Found user - %v", user)
	}
	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting user resources")
	for i := range s.users {
		err := s.scope.KeycloakClient.DeleteUser(ctx, s.scope.Token(), s.scope.RealmName(), *s.users[i].ID)
		if err != nil {
			log.Info("Unable to delete user - %s", err)
		}
		log.Info("Deleted user - %s", s.users[i].Username)
	}

	return nil
}
