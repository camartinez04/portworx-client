package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	//Creates the channel MailChan from the model MailData
}

var (
	useTls = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	token  = flag.String("token", "", "Authorization token if any")
)

var brokerURL = os.Getenv("BROKER_URL")
