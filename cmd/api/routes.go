package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.ServeFiles("/v1/image/*filepath", http.Dir("internal/images"))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/grenades/", app.getAllGrenadesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/grenades/:id", app.getGrenadeHandler)
	router.HandlerFunc(http.MethodPost, "/v1/grenades", app.createGrenadeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/grenades/:id", app.updateGrenadeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/grenades/:id", app.deleteGrenadeHandler)

	router.HandlerFunc(http.MethodPost, "/v1/grenades/:id/images", app.uploadImageHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/images/:id", app.deleteImageHandler)
	
	return router
}
