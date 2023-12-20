package realms

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling realm resources")
	realm, err := s.scope.KeycloakClient.GetRealm(ctx, s.scope.KeycloakToken.AccessToken, s.scope.RealmName())
	if err != nil {
		return err
	}
	log.Info("Realm Name - %s", realm.DisplayName)
	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting realm resources not yet supported")

	return nil
}
