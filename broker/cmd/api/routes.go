package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the routes for the API
func (app *AppConfig) routes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/getpxclustercapacity", app.getPXClusterCapacityHTTP)
	mux.Get("/getpxcluster", app.getPXClusterHTTP)
	mux.Get("/getpxclusteralarms", app.getPXClusterAlarmsHTTP)

	mux.Post("/postcreatevolume", app.postCreateNewVolumeHTTP)
	mux.Get("/getvolumeinfo/{volume_id}", app.getVolumeInfoHTTP)
	mux.Patch("/patchvolumesize/{volume_id}", app.patchUpdateVolumeSizeHTTP)
	mux.Patch("/patchvolumeioprofile/{volume_id}", app.patchUpdateVolumeIOProfileHTTP)
	mux.Patch("/patchvolumehalevel/{volume_id}", app.patchUpdateVolumeHALevelHTTP)
	mux.Patch("/patchvolumesharedv4/{volume_id}", app.patchUpdateVolumeSharedv4HTTP)
	mux.Patch("/patchvolumesharedv4service/{volume_id}", app.patchUpdateVolumeSharedvService4HTTP)
	mux.Patch("/patchvolumenodiscard/{volume_id}", app.patchUpdateVolumeNoDiscardHTTP)
	mux.Delete("/deletevolume/{volume_id}", app.deleteVolumeHTTP)

	mux.Get("/getvolumeid/{volume_name}", app.getVolumeIDsHTTP)
	mux.Get("/getnodesofvolume/{volume_name}", app.getNodesOfVolumeHTTP)
	mux.Get("/getinspectvolume/{volume_name}", app.getInspectVolumeHTTP)
	mux.Get("/getvolumeusage/{volume_name}", app.getVolumeUsageHTTP)

	mux.Get("/getallvolumesinfo", app.getAllVolumesInfoHTTP)
	mux.Get("/getallvolumes", app.getAllVolumesHTTP)
	mux.Get("/getallvolumescomplete", app.getAllVolumesCompleteHTTP)
	mux.Get("/getreplicaspernode/{node_id}", app.getReplicasPerNodeHTTP)

	mux.Get("/getlistofnodes", app.getListOfNodesHTTP)
	mux.Get("/getallnodesinfo", app.getAllNodesInfoHTTP)
	mux.Get("/getnodeinfo/{node_id}", app.getNodeInfoHTTP)
	mux.Get("/getcloudsnaps/{volume_id}", app.getCloudSnapsHTTP)
	mux.Get("/getallcloudsnaps", app.getAllCloudSnapsHTTP)
	mux.Get("/getinspectawscloudcreds", app.getInspectAWSCloudCredentialHTTP)
	mux.Post("/postcreateawscloudcreds", app.postCreateAWSCloudCredentialHTTP)
	mux.Delete("/deleteawscloudcreds", app.deleteAWSCloudCredentialHTTP)
	mux.Post("/postcreatecloudsnap", app.postCreateCloudSnapHTTP)
	mux.Post("/postcreatelocalsnap", app.postCreateLocalSnapHTTP)

	return mux

}
