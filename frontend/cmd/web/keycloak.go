package main

import (
	"github.com/Nerzal/gocloak/v7"
)

func newKeycloak() *keycloak {
	return &keycloak{
		gocloak:      gocloak.NewClient(KeycloakURL),
		clientId:     KeycloakClientID,
		clientSecret: KeycloakSecret,
		realm:        KeycloakRealm,
	}
}
