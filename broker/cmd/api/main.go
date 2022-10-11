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

// GetRequestMetadata gets the current request metadata.
func (t OpenStorageSdkToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "bearer " + *token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (t OpenStorageSdkToken) RequireTransportSecurity() bool {
	return *useTls
}

// main is the entry point for the API
func main() {

	flag.Parse()

	contextToken := OpenStorageSdkToken{}

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gob.Register(loginRequest{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.InProduction = false

	app.Session = session

	if *useTls {
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

	if len(*token) != 0 {
		// Add token interceptor
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(contextToken))
	}

	conn, err := grpc.Dial(*address, dialOptions...)
	if err != nil {
		log.Panicf("Error trying to establish gRPC connection to address %s: %v", err, address)
		os.Exit(1)
	}

	app.NewKeycloak = newKeycloak()

	// Sends the grpc connection that we've created to AppConfig
	app := AppConfig{
		Conn:        conn,
		Session:     session,
		NewKeycloak: newKeycloak(),
	}

	NewHandlers(&app)

	log.Printf("Connected to Portworx's OpenStorage via gRPC to %s", *address)

	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	log.Println("Starting Broker server on port", webPort)

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
