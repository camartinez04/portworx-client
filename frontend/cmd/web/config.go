package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/alexedwards/scs/v2"
)

const PortNumber = ":8082"

var App AppConfig
var InfoLog *log.Logger
var ErrorLog *log.Logger
var Repo *Repository
var Session *scs.SessionManager
var KeycloakToken string
var KeycloakRefreshToken string

var BrokerURL = os.Getenv("BROKER_URL")
var KeycloakURL = os.Getenv("KEYCLOAK_URL")
var KeycloakClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
var KeycloakSecret = os.Getenv("KEYCLOAK_SECRET")
var KeycloakRealm = os.Getenv("KEYCLOAK_REALM")

var PathToTemplates = "./static/templates"

var (
	UseTls = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	//token  = flag.String("token", "", "Authorization token if any")
)

// Map of functions available to the templates
var Functions = template.FuncMap{
	"humanDate":        HumanDate,
	"formatDate":       FormatDate,
	"iterate":          Iterate,
	"add":              Add,
	"divide":           Divide,
	"resizeVolume":     ResizeVolume,
	"split":            strings.Split,
	"removeDuplicates": RemoveDuplicateStr,
	"dateFromUnix":     DateFormat,
	"byteToGigabyte":   BytesToGB,
	"byteToMegabyte":   BytesToMB,
}
