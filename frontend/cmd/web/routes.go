package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the routing for the application.
func routes(app *AppConfig) http.Handler {

	mux := chi.NewRouter()

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

	mux.Route("/frontend", func(mux chi.Router) {

		mux.Get("/", Repo.Cluster)

		mux.Get("/oauth/callback", Repo.Cluster)

		mux.Get("/cluster", Repo.Cluster)

		mux.Get("/volumes", Repo.Volumes)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/volume/{volume_id}", Repo.VolumeInformation)

		//mux.Post("/volume/{volume_id}", Repo.VolumeInformation)

		//mux.Delete("/volume/{volume_id}", Repo.VolumeInformation)

		mux.Get("/nodes", Repo.Nodes)

		middleware.SetHeader("Node-ID", "{node_id}")
		mux.Get("/node/{node_id}", Repo.NodeInformation)

		mux.Get("/snapshots", Repo.GetAllSnaps)

		mux.Get("/snapshot/{bucket}/{snap_id}", Repo.SpecificSpapInformation)

		mux.Get("/cloud-credentials", Repo.CloudCredentials)

		mux.Get("/cloud-credential/{cloud_cred_id}", Repo.CloudCredentialsInformation)

		mux.Get("/create-credentials", Repo.CreateCloudCredentials)

		mux.Post("/create-credentials", Repo.PostCreateCloudCredentials)

		mux.Get("/create-volume", Repo.CreateVolume)

		mux.Post("/create-volume", Repo.PostCreateVolume)

		mux.Get("/documentation", Repo.Documentation)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/delete-volume/{volume_id}", Repo.DeleteVolume)

		mux.Get("/update-volume-halevel/{volume_id}/{ha-level}", Repo.UpdateVolumeHALevelHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/update-volume-size/{volume_id}/{size}", Repo.UpdateVolumeSizeHTTP)

		middleware.SetHeader("Volume-ID", "{volume_id}")
		mux.Get("/update-volume-ioprofile/{volume_id}/{ioprofile}", Repo.UpdateVolumeIOProfileHTTP)

		fileServer := http.FileServer(http.Dir("./static/"))

		mux.Handle("/static/*", http.StripPrefix("/frontend/static", fileServer))

	})

	return mux
}
