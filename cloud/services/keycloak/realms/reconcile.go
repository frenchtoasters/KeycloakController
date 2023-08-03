package realms

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling realm resources")
	_, _, err := s.realms.Get(ctx, s.scope.RealmName())
	if err != nil {
		log.Info("Realm does not yet exist")
	}
	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting realm resources not yet supported")

	return nil
}
