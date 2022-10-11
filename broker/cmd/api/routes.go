package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the routes for the API
func (App *AppConfig) routes() http.Handler {

	keycloak := App.NewKeycloak

	mux := chi.NewRouter()

	mdw := newMiddleware(keycloak)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "ACCEPTED"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Route("/", func(mux chi.Router) {

		mux.Post("/login", App.postLoginHTTP)
		mux.Get("/login", App.postLoginHTTP)
		mux.Get("/logout", App.getLogoutHTTP)
	})

	mux.Route("/broker", func(mux chi.Router) {

		mux.Use(mdw.AuthKeycloak)

		mux.Get("/getpxclustercapacity", App.getPXClusterCapacityHTTP)
		mux.Get("/getpxcluster", App.getPXClusterHTTP)
		mux.Get("/getpxclusteralarms", App.getPXClusterAlarmsHTTP)

		mux.Post("/postcreatevolume", App.postCreateNewVolumeHTTP)
		mux.Get("/getvolumeinfo/{volume_id}", App.getVolumeInfoHTTP)
		mux.Patch("/patchvolumesize/{volume_id}", App.patchUpdateVolumeSizeHTTP)
		mux.Patch("/patchvolumeioprofile/{volume_id}", App.patchUpdateVolumeIOProfileHTTP)
		mux.Patch("/patchvolumehalevel/{volume_id}", App.patchUpdateVolumeHALevelHTTP)
		mux.Patch("/patchvolumesharedv4/{volume_id}", App.patchUpdateVolumeSharedv4HTTP)
		mux.Patch("/patchvolumesharedv4service/{volume_id}", App.patchUpdateVolumeSharedvService4HTTP)
		mux.Patch("/patchvolumenodiscard/{volume_id}", App.patchUpdateVolumeNoDiscardHTTP)
		mux.Delete("/deletevolume/{volume_id}", App.deleteVolumeHTTP)

		mux.Get("/getvolumeid/{volume_name}", App.getVolumeIDsHTTP)
		mux.Get("/getnodesofvolume/{volume_name}", App.getNodesOfVolumeHTTP)
		mux.Get("/getinspectvolume/{volume_name}", App.getInspectVolumeHTTP)
		mux.Get("/getvolumeusage/{volume_name}", App.getVolumeUsageHTTP)

		mux.Get("/getallvolumesinfo", App.getAllVolumesInfoHTTP)
		mux.Get("/getallvolumes", App.getAllVolumesHTTP)
		mux.Get("/getallvolumescomplete", App.getAllVolumesCompleteHTTP)
		mux.Get("/getreplicaspernode/{node_id}", App.getReplicasPerNodeHTTP)

		mux.Get("/getlistofnodes", App.getListOfNodesHTTP)
		mux.Get("/getallnodesinfo", App.getAllNodesInfoHTTP)
		mux.Get("/getnodeinfo/{node_id}", App.getNodeInfoHTTP)
		mux.Get("/getcloudsnaps/{volume_id}", App.getCloudSnapsHTTP)
		mux.Get("/getallcloudsnaps", App.getAllCloudSnapsHTTP)
		mux.Get("/getinspectawscloudcreds", App.getInspectAWSCloudCredentialHTTP)
		mux.Post("/postcreateawscloudcreds", App.postCreateAWSCloudCredentialHTTP)
		mux.Delete("/deleteawscloudcreds", App.deleteAWSCloudCredentialHTTP)
		mux.Post("/postcreatecloudsnap", App.postCreateCloudSnapHTTP)
		mux.Post("/postcreatelocalsnap", App.postCreateLocalSnapHTTP)
		mux.Get("/getallcloudcredsids", App.getAllCloudCredentialIDsHTTP)
		mux.Get("/getspecificcloudsnapshot", App.getSpecificCloudSnapshotHTTP)
		mux.Delete("/deletecloudsnap", App.deleteCloudSnapHTTP)

	})

	return mux

}
