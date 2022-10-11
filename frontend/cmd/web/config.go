package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8081"

var app AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
var Repo *Repository
var session *scs.SessionManager
var keycloakToken string
var keycloakRefreshToken string

var brokerURL = os.Getenv("BROKER_URL")
var KeycloakURL = os.Getenv("KEYCLOAK_URL")
var KeycloakClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
var KeycloakSecret = os.Getenv("KEYCLOAK_SECRET")
var KeycloakRealm = os.Getenv("KEYCLOAK_REALM")

var pathToTemplates = "./static/templates"

var (
	useTls = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	//token  = flag.String("token", "", "Authorization token if any")
)

// Map of functions available to the templates
var functions = template.FuncMap{
	"humanDate":        HumanDate,
	"formatDate":       FormatDate,
	"iterate":          Iterate,
	"add":              Add,
	"divide":           Divide,
	"resizeVolume":     ResizeVolume,
	"split":            strings.Split,
	"removeDuplicates": RemoveDuplicateStr,
	"dateFromUnix":     DateFormat,
}
