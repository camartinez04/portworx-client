package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *AppConfig) http.Handler {

	//mux := pat.New()
	//mux.Get("/hello-world", http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Route("/frontend", func(mux chi.Router) {

		mux.Get("/cluster", Repo.Cluster)

		mux.Get("/volumes", Repo.Volumes)

		mux.Get("/volume/{volume_id}", Repo.VolumeInformation)

		mux.Post("/volume/{volume_id}", Repo.VolumeInformation)

		mux.Get("/nodes", Repo.Nodes)

		mux.Get("/node/{node_id}", Repo.NodeInformation)

		mux.Get("/snapshots", Repo.Snaps)

		mux.Get("/snapshot/{snap_name}", Repo.SnapsInformation)

		mux.Get("/storage-pools", Repo.StoragePools)

		mux.Get("/stogage-pool/{stg_name}", Repo.StoragePoolsInformation)

		mux.Get("/create-volume", Repo.GetCreateVolume)

		mux.Post("/create-volume", Repo.PostCreateVolume)

		fileServer := http.FileServer(http.Dir("./static/"))

		mux.Handle("/static/*", http.StripPrefix("/frontend/static", fileServer))

	})

	return mux
}
