package main

import (
	"context"
	"crypto/x509"
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GetRequestMetadata gets the current request metadata. Ensure not adding /n to the end of the token, otherwise it will fail
func (t OpenStorageSdkToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "bearer " + t.Token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (t OpenStorageSdkToken) RequireTransportSecurity() bool {
	return *UseTls
}

// main is the entry point for the API
func main() {

	flag.Parse()

	contextToken := OpenStorageSdkToken{}

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gob.Register(LoginRequest{})

	Session = scs.New()
	Session.Lifetime = 24 * time.Hour
	Session.Cookie.Persist = true
	Session.Cookie.SameSite = http.SameSiteLaxMode
	Session.Cookie.Secure = App.InProduction

	App.InProduction = false

	App.Session = Session

	if *UseTls {
		// Setup a connection
		capool, err := x509.SystemCertPool()
		if err != nil {
			log.Panicf("Failed to load system certs: %v\n", err)
			os.Exit(1)
		}
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(capool, ""),
		)}
	}

	if len(*Token) != 0 {
		// Set token
		contextToken.Token = *Token

		// Add token to dial options
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(contextToken))
	}

	conn, err := grpc.Dial(*Address, dialOptions...)
	if err != nil {
		log.Panicf("Error trying to establish gRPC connection to address %s: %v", err, Address)
		os.Exit(1)
	}

	App.NewKeycloak = NewKeycloak()

	// Sends the grpc connection that we've created to AppConfig
	app := AppConfig{
		Conn:        conn,
		Session:     Session,
		NewKeycloak: NewKeycloak(),
	}

	NewHandlers(&app)

	log.Printf("Connected to Portworx's OpenStorage via gRPC to %s", *Address)

	srv := &http.Server{
		Addr:    WebPort,
		Handler: app.routes(),
	}

	log.Println("Starting Broker server on port", WebPort)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Panic(err)
		}
	}()

}
