package realms

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling realm resources")
	log.Info("Realm Name - %s", s.scope.Realms())
	return nil

}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting realm resources not yet supported")

	return nil
}
