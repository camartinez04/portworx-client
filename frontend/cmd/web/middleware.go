package main

import (
	"log"
	"net/http"

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
		Secure:   App.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "CSRF token missing or invalid", http.StatusBadRequest)
	}))

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return Session.LoadAndSave(next)
}

// NewMiddleware creates a new middleware with Keycloak
func NewMiddleware(keycloak *Keycloak) *KeyCloakMiddleware {

	return &KeyCloakMiddleware{keycloak: keycloak}
}
