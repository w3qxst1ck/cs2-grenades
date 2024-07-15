package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/w3qxst1ck/cs2-grenades/internal/data"
	"github.com/w3qxst1ck/cs2-grenades/internal/validator"
)

func (app *application) getGrenadeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// get grenade
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

	// get images for grenade
	images, err := app.models.Images.GetByGrenadeID(grenade.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create Image.ImageURL for images
	app.createImagesURL(images)

	grenade.Images = images

	app.cache.Set(r.URL.Path, envelope{"grenade": grenade}, 0)

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

func (app *application) updateGrenadeHandler(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Map         *string `json:"map"`
		Title       *string `json:"title"`
		Type        *string `json:"type"`
		Side        *string `json:"side"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Map != nil {
		grenade.Map = *input.Map
	}

	if input.Title != nil {
		grenade.Title = *input.Title
	}

	if input.Type != nil {
		grenade.Type = *input.Type
	}

	if input.Description != nil {
		grenade.Description = *input.Description
	}

	if input.Side != nil {
		grenade.Side = *input.Side
	}

	v := validator.New()
	if data.ValidateGrenade(grenade, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Erorrs)
		return
	}

	err = app.models.Grenades.Update(grenade)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
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

func (app *application) deleteGrenadeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	images, err := app.models.Images.GetByGrenadeID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err = app.deleteGrenadeImages(images); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Grenades.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "grenade successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllGrenadesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Map  string
		Side string
		Type string
		data.Filters
	}

	qs := r.URL.Query()
	input.Map = app.readString(qs, "map", "")
	input.Side = app.readString(qs, "side", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "map", "side", "type", "-id"}

	v := validator.New()
	v.Check(v.In(input.Filters.Sort, input.Filters.SortSafeList), "sort", "invalid sort value")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Erorrs)
		return
	}

	grenades, err := app.models.Grenades.GetAll(input.Map, input.Side, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var wg sync.WaitGroup

	for i := range grenades {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			grenade := grenades[i]
			images, err := app.models.Images.GetByGrenadeID(grenade.ID)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			app.createImagesURL(images)
			grenade.Images = images
		}(i)
	}

	wg.Wait()

	cachePath := r.URL.Path + qs.Encode()
	app.cache.Set(cachePath, envelope{"grenades": grenades}, 0)

	err = app.writeJSON(w, http.StatusOK, envelope{"grenades": grenades}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createImagesURL(images []*data.Image) {
	for i := range images {
		images[i].ImageURL = fmt.Sprintf("%s%s", app.config.storageS3.DownloadUrl, images[i].Name)
	}
}

func (app *application) deleteGrenadeImages(images []*data.Image) error {
	for _, image := range images {
		err := os.Remove(fmt.Sprintf("%s%s", app.config.imagesDir, image.Name))
		if err != nil {
			return err
		}
	}
	return nil
}
