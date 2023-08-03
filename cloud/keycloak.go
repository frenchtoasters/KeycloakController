package cloud

import (
	"context"

	gokeycloak "github.com/Nerzal/gocloak/v13"
)

type goKeycloakInterface interface {
	Realms() Realms
}

type Realms interface {
	Get(ctx context.Context, token string) (int, []*gokeycloak.RealmRepresentation, error)
}
