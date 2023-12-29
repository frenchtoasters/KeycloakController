package realms

import (
	"context"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/utils"
	gokeycloak "github.com/Nerzal/gocloak/v13"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling realm resources")
	realm, err := s.scope.KeycloakClient.GetRealm(ctx, s.scope.Token(), s.scope.RealmName())
	if err != nil {
		if err.(*gokeycloak.APIError).Code == 404 {
			log.Info("Realm not found creating new realm")
			realmRep := utils.DefaultRealmRep(&s.scope.Keycloak.Spec.RealmName)
			realm, err := s.scope.KeycloakClient.CreateRealm(ctx, s.scope.Token(), realmRep)
			if err != nil {
				return err
			}
			log.Info("Realm Created", "RealmName", realm)
			return nil
		} else {
			return err
		}
	}
	s.scope.Keycloak.Status.RealmName = *realm.Realm
	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting realm")
	err := s.scope.KeycloakClient.DeleteRealm(ctx, s.scope.Token(), s.scope.RealmName())
	if err != nil {
		return err
	}

	return nil
}
