package main

import (
	"net/http"
)

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// newMiddleware creates a new middleware with Keycloak
func newMiddleware(keycloak *keycloak) *keyCloakMiddleware {

	return &keyCloakMiddleware{keycloak: keycloak}
}
