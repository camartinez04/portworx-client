package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
)

// WriteToConsole writes something to the console
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protection
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "CSRF token missing or invalid", http.StatusBadRequest)
	}))

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// newMiddleware creates a new middleware with Keycloak
func newMiddleware(keycloak *keycloak) *keyCloakMiddleware {

	return &keyCloakMiddleware{keycloak: keycloak}
}

// extractBearerToken extracts the Bearer token from the Authorization header
func (auth *keyCloakMiddleware) extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

// AuthKeycloak is a middleware to check if the user is authenticated and check the JWT token
func (auth *keyCloakMiddleware) AuthKeycloak(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		// Check if the user is authenticated
		if keycloakToken == "" {
			session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "/portworx/login", http.StatusSeeOther)
			return
		}

		r.Header.Add("Authorization", "Bearer "+keycloakToken)

		token := r.Header.Get("Authorization")

		// Extract Bearer token
		token = auth.extractBearerToken(token)

		if token == "" {
			session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "/portworx/login", http.StatusSeeOther)
			return
		}

		// Call Keycloak API to verify the access token
		result, err := auth.keycloak.gocloak.RetrospectToken(context.Background(), token, auth.keycloak.clientId, auth.keycloak.clientSecret, auth.keycloak.realm)
		if err != nil {
			session.Put(r.Context(), "error", fmt.Sprintf("Invalid or malformed token: %s", err.Error()))
			http.Redirect(w, r, "/portworx/login", http.StatusSeeOther)
			return
		}

		// Decode the token and validate it
		_, _, err = auth.keycloak.gocloak.DecodeAccessToken(context.Background(), token, auth.keycloak.realm)
		if err != nil {
			session.Put(r.Context(), "error", fmt.Sprintf("Invalid or malformed Token when decoding it %s", err.Error()))
			http.Redirect(w, r, "/portworx/login", http.StatusSeeOther)
			return
		}

		// Check if the token isn't expired and valid
		if !*result.Active {
			session.Put(r.Context(), "error", "Invalid or expired Token")
			http.Redirect(w, r, "/portworx/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

// newController creates a new controller
func newController(keycloak *keycloak) *controller {
	return &controller{
		keycloak: keycloak,
	}
}
