package main

import (
	"context"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (t OpenStorageSdkToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "bearer " + *token,
	}, nil
}

func (t OpenStorageSdkToken) RequireTransportSecurity() bool {
	return *useTls
}

func main() {

	flag.Parse()

	contextToken := OpenStorageSdkToken{}

	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	if *useTls {
		// Setup a connection
		capool, err := x509.SystemCertPool()
		if err != nil {
			fmt.Printf("Failed to load system certs: %v\n")
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
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Sends the grpc connection that we've created to AppConfig
	app := AppConfig{
		Conn: conn,
	}

	log.Println("Starting server on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Panic(err)
		}
	}()

}
