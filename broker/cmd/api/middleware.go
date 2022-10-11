package main

import (
	"net/http"
)

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return Session.LoadAndSave(next)
}

// newMiddleware creates a new middleware with Keycloak
func newMiddleware(keycloak *Keycloak) *KeyCloakMiddleware {

	return &KeyCloakMiddleware{keycloak: keycloak}
}
