package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v11"
)

// NewKeycloak creates a new keycloak call
func NewKeycloak() *Keycloak {
	return &Keycloak{
		gocloak:      gocloak.NewClient(KeycloakURL),
		clientId:     KeycloakClientID,
		clientSecret: KeycloakSecret,
		realm:        KeycloakRealm,
	}
}

// ExtractBearerToken extracts the Bearer token from the Authorization header
func (auth *KeyCloakMiddleware) ExtractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

// AuthKeycloak is a middleware to check if the user is authenticated and check the JWT token
func (auth *KeyCloakMiddleware) AuthKeycloak(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		// Check if the user is authenticated
		if KeycloakToken == "" {
			Session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		r.Header.Add("Authorization", "Bearer "+KeycloakToken)

		token := r.Header.Get("Authorization")

		// Extract Bearer token
		token = auth.ExtractBearerToken(token)

		if token == "" {
			Session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
			return
		}

		// Call Keycloak API to verify the access token
		result, err := auth.keycloak.gocloak.RetrospectToken(context.Background(), token, auth.keycloak.clientId, auth.keycloak.clientSecret, auth.keycloak.realm)
		if err != nil {
			Session.Put(r.Context(), "error", fmt.Sprintf("Invalid or malformed token: %s", err.Error()))
			http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
			return
		}

		// Decode the token and validate it
		_, _, err = auth.keycloak.gocloak.DecodeAccessToken(context.Background(), token, auth.keycloak.realm)
		if err != nil {
			Session.Put(r.Context(), "error", fmt.Sprintf("Invalid or malformed Token when decoding it %s", err.Error()))
			http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
			return
		}

		// Check if the token isn't expired and valid
		if !*result.Active {
			Session.Put(r.Context(), "error", "Invalid or expired Token")
			http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

// NewController creates a new controller
func NewController(keycloak *Keycloak) *Controller {
	return &Controller{
		keycloak: keycloak,
	}
}
