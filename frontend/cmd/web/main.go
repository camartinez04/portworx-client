package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8081"

var app AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
var session *scs.SessionManager
var keycloakToken string
var keycloakRefreshToken string

func main() {

	//get the template cache from the app config
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	gob.Register(map[string]int{})
	gob.Register(CreateVolume{})
	gob.Register(CreateCloudCredentials{})
	gob.Register(CreateCloudSnap{})
	gob.Register(loginRequest{})

	app.NewKeycloak = newKeycloak()

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.InProduction = false

	app.Session = session

	app.TemplateCache = tc
	app.UseCache = false

	repo := NewRepo(&app)

	NewHandlers(repo)
	NewRenderer(&app)
	NewHelpers(&app)

	log.Println("Starting Frontend on port", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
