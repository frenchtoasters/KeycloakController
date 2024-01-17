package realmroles

import (
	"context"
	"fmt"

	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/utils"
	"github.com/Nerzal/gocloak/v13"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling realm role resources")
	compositeRoleNames := []string{"manage-users", "view-users", "view-clients", "view-events", "view-identity-providers"}

	adminRoles, err := s.scope.KeycloakClient.GetCompositeRealmRoles(ctx, s.scope.Token(), "master", "admin")
	if err != nil {
		return fmt.Errorf("error getting composite realm roles: %s", err)
	}

	clientParams := &gocloak.GetClientsParams{
		Search:   gocloak.BoolP(true),
		ClientID: gocloak.StringP("realm-management"),
	}
	clientId, err := s.scope.KeycloakClient.GetClients(ctx, s.scope.Token(), s.scope.RealmName(), *clientParams)
	if err != nil {
		return fmt.Errorf("error getting clients: %s", err)
	}

	roleComposites := []gocloak.Role{}
	for i := range adminRoles {
		if slices.Contains(compositeRoleNames, *adminRoles[i].Name) {
			role, err := s.scope.KeycloakClient.GetClientRole(ctx, s.scope.Token(), s.scope.RealmName(), *clientId[0].ID, *adminRoles[i].Name)
			if err != nil {
				return fmt.Errorf("error getting client role: %s", err)
			}
			roleComposites = append(roleComposites, *role)
		}
	}

	tenantAdminRole := gocloak.Role{
		Name:      gocloak.StringP("tenant-realm-admin"),
		Composite: gocloak.BoolP(true),
	}

	adminRoleParams := gocloak.GetRoleParams{
		Search: gocloak.StringP("tenant-realm-admin"),
	}

	// Check if role exists or not
	adminRole, err := s.scope.KeycloakClient.GetRealmRoles(ctx, s.scope.Token(), s.scope.RealmName(), adminRoleParams)
	if err != nil {
		return fmt.Errorf("error when checking for realm role: %s", err)
	}

	// Create role if no roles by that name are found
	if len(adminRole) == 0 {
		adminRoleName, err := s.scope.KeycloakClient.CreateRealmRole(ctx, s.scope.Token(), s.scope.RealmName(), tenantAdminRole)
		if err != nil {
			return fmt.Errorf("error creating realm role: %s", err)
		}
		log.Info(fmt.Sprintf("Admin Role Created: %s", adminRoleName))
	}
	log.Info(fmt.Sprintf("Reconciling admin role composites"))

	adminRoleId, err := s.scope.KeycloakClient.GetRealmRoles(ctx, s.scope.Token(), s.scope.RealmName(), adminRoleParams)
	if err != nil {
		return fmt.Errorf("error getting admin role: %s", err)
	}

	if err := s.scope.KeycloakClient.AddClientRoleComposite(ctx, s.scope.Token(), s.scope.RealmName(), *adminRoleId[0].ID, roleComposites); err != nil {
		return fmt.Errorf("error adding client role composite: %s", err)
	}

	s.scope.Keycloak.Spec.AdminRole = utils.RoleTransform(adminRoleId[0])
	log.Info("Admin role composites added")

	return nil
}

func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting realm roles")
	return nil
}
