package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/w3qxst1ck/cs2-grenades/internal/data"
	"github.com/w3qxst1ck/cs2-grenades/internal/validator"
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

func (app *application) createGrenadeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Map         string `json:"map"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Side        string `json:"side"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	grenade := &data.Grenade{
		Map:         input.Map,
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Side:        input.Side,
	}

	v := validator.New()
	data.ValidateGrenade(grenade, v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Erorrs)
		return
	}

	err = app.models.Grenades.Insert(grenade)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/grenades/%d", grenade.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"grenade": grenade}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
