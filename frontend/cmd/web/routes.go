package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the routing for the application.
func routes(app *AppConfig) http.Handler {

	keycloak := Repo.App.NewKeycloak

	mux := chi.NewRouter()

	mdw := newMiddleware(keycloak)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Route("/", func(mux chi.Router) {

		mux.Get("/", Repo.GetLoginHTTP)

		mux.Post("/", Repo.PostLoginHTTP)
	})

	mux.Route("/portworx", func(mux chi.Router) {

		mux.Get("/", Repo.GetLoginHTTP)

		mux.Post("/", Repo.PostLoginHTTP)

		mux.Get("/login", Repo.GetLoginHTTP)

		mux.Post("/login", Repo.PostLoginHTTP)

		fileServer := http.FileServer(http.Dir("./static/"))

		mux.Handle("/static/*", http.StripPrefix("/portworx/static", fileServer))

	})

	mux.Route("/portworx/client", func(mux chi.Router) {

		mux.Use(mdw.AuthKeycloak)

		mux.Get("/", Repo.ClusterHTTP)

		mux.Post("/logout", Repo.LogoutHTTP)

		mux.Get("/oauth/callback", Repo.ClusterHTTP)

		mux.Get("/cluster", Repo.ClusterHTTP)

		mux.Get("/volumes", Repo.VolumesHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/volume/{volume_id}", Repo.VolumeInformationHTTP)

		mux.Get("/nodes", Repo.NodesHTTP)

		middleware.SetHeader("Node-ID", "{node_id}")
		mux.Get("/node/{node_id}", Repo.NodeInformationHTTP)

		mux.Get("/snapshots", Repo.GetAllSnapsHTTP)

		mux.Get("/snapshot/{bucket}/{snap_id}", Repo.SpecificSpapInformationHTTP)

		mux.Get("/cloud-credentials", Repo.CloudCredentialsHTTP)

		mux.Get("/cloud-credential/{cloud_cred_id}", Repo.CloudCredentialsInformationHTTP)

		mux.Get("/create-credentials", Repo.CreateCloudCredentialsHTTP)

		mux.Post("/create-credentials", Repo.PostCreateCloudCredentialsHTTP)

		mux.Get("/create-volume", Repo.CreateVolumeHTTP)

		mux.Post("/create-volume", Repo.PostCreateVolumeHTTP)

		mux.Get("/create-cloudsnap", Repo.CreateCloudSnapHTTP)

		mux.Post("/create-cloudsnap", Repo.PostCreateCloudSnapHTTP)

		mux.Get("/documentation", Repo.DocumentationHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/delete-volume/{volume_id}", Repo.DeleteVolumeHTTP)

		mux.Get("/update-volume-halevel/{volume_id}/{ha-level}", Repo.UpdateVolumeHALevelHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/update-volume-size/{volume_id}/{size}", Repo.UpdateVolumeSizeHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/update-volume-ioprofile/{volume_id}/{ioprofile}", Repo.UpdateVolumeIOProfileHTTP)

		fileServer := http.FileServer(http.Dir("./static/"))

		mux.Handle("/static/*", http.StripPrefix("/portworx/client/static", fileServer))

	})

	return mux
}
