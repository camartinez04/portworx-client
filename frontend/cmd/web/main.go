package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

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
	gob.Register(LoginRequest{})

	App.NewKeycloak = NewKeycloak()

	Session = scs.New()
	Session.Lifetime = 24 * time.Hour
	Session.Cookie.Persist = true
	Session.Cookie.SameSite = http.SameSiteLaxMode
	Session.Cookie.Secure = App.InProduction

	App.InProduction = false

	App.Session = Session

	App.TemplateCache = tc
	App.UseCache = false

	repo := NewRepo(&App)

	NewHandlers(repo)
	NewRenderer(&App)
	NewHelpers(&App)

	log.Println("Starting Frontend on port", PortNumber)

	srv := &http.Server{
		Addr:    PortNumber,
		Handler: routes(&App),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
