package cloud

import (
	"context"

	gokeycloak "github.com/Nerzal/gocloak/v13"
)

type goKeycloakInterface interface {
	Realms() Realms
	Users() Users
}

type Realms interface {
	Get(ctx context.Context, token string) (int, []*gokeycloak.RealmRepresentation, error)
}

type Users interface {
	Get(ctx context.Context, token string) (int, []*gokeycloak.User, error)
}

type Groups interface {
	Get(ctx context.Context, token string) (int, []*gokeycloak.Group, error)
}

type RealmRoles interface {
	Get(ctx context.Context, token string) (int, []*gokeycloak.Role, error)
}
