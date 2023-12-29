package users

import (
	"context"
	"fmt"

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
		if len(user) != 1 {
			log.Info("Reconciling users")
			userId, err := s.scope.KeycloakClient.CreateUser(ctx, s.scope.Token(), s.scope.RealmName(), utils.UserTransform(s.users[i]))
			if err != nil {
				log.Info("Unable to create user", "error", err)
				return err
			}
			log.Info(fmt.Sprintf("Created user - [%s]", userId))
			continue
		}
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Found user - [%s]", *user[0].Username))
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
