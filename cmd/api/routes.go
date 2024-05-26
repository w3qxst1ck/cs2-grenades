package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.ServeFiles("/v1/image/*filepath", http.Dir("internal/images"))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/grenades/", app.checkCache(app.getAllGrenadesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/grenades/:id", app.checkCache(app.getGrenadeHandler))
	router.HandlerFunc(http.MethodPost, "/v1/grenades", app.createGrenadeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/grenades/:id", app.updateGrenadeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/grenades/:id", app.deleteGrenadeHandler)

	router.HandlerFunc(http.MethodPost, "/v1/grenades/:id/images", app.uploadImageHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/images/:id", app.deleteImageHandler)

	return app.recoverPanic(app.rateLimit((router)))
}
