package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *AppConfig) routes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/getclustercapacity", app.getClusterCapacityHTTP)
	mux.Get("/getclusteruuid", app.getClusterUUIDHTTP)
	mux.Get("/getvolumeid/{volume_name}", app.getVolumeIDsHTTP)
	mux.Get("/getnodesofvolume/{volume_name}", app.getNodesOfVolumeHTTP)
	mux.Get("/getinspectvolume/{volume_name}", app.getInspectVolumeHTTP)
	mux.Get("/getvolumeusage/{volume_name}", app.getVolumeUsageHTTP)
	mux.Get("/getlistofnodes", app.getListOfNodesHTTP)
	mux.Get("/getreplicaspernode/{node_id}", app.getReplicasPerNodeHTTP)
	mux.Get("/getallvolumes", app.getAllVolumesHTTP)
	mux.Get("/getallvolumescomplete", app.getAllVolumesCompleteHTTP)
	mux.Get("/getallnodesinfo", app.getAllNodesInfoHTTP)
	mux.Get("/getnodeinfo/{node_id}", app.getNodeInfoHTTP)
	mux.Get("/getallvolumesinfo", app.getAllVolumesInfoHTTP)
	mux.Get("/getvolumeinfo/{node_id}", app.getVolumeInfoHTTP)
	mux.Post("/postcreatevolume", app.postCreateNewVolumeHTTP)

	return mux

}
