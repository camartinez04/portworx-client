package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8081"

var app AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	//get the template cache from the app config
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := NewRepo(&app)

	NewHandlers(repo)
	NewRenderer(&app)
	NewHelpers(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
