package main

import (
	"flag"
	"html/template"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v7"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	InProduction  bool
	NewKeycloak   *keycloak
	//Creates the channel MailChan from the model MailData
}

var KeycloakURL = os.Getenv("KEYCLOAK_URL")
var KeycloakClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
var KeycloakSecret = os.Getenv("KEYCLOAK_SECRET")
var KeycloakRealm = os.Getenv("KEYCLOAK_REALM")

var (
	useTls = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	//token  = flag.String("token", "", "Authorization token if any")
)

var brokerURL = os.Getenv("BROKER_URL")

type keycloak struct {
	gocloak      gocloak.GoCloak // keycloak client
	clientId     string          // clientId specified in Keycloak
	clientSecret string          // client secret specified in Keycloak
	realm        string          // realm specified in Keycloak
}

type keyCloakMiddleware struct {
	keycloak *keycloak
	Session  *scs.SessionManager
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type controller struct {
	keycloak *keycloak
}
