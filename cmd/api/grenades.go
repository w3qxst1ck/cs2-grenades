package main

import (
	"errors"
	"net/http"

	"github.com/w3qxst1ck/cs2-grenades/internal/data"
)

func (app *application) getGrenadeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	grenade, err := app.models.Grenades.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"grenade": grenade}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
