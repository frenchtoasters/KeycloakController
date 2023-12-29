package utils

import "github.com/Nerzal/gocloak/v13"

func DefaultRealmRep(name *string) gocloak.RealmRepresentation {
	realm := &gocloak.RealmRepresentation{
		Realm: name,
	}
	return *realm
}
