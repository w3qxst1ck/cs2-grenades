package main

import (
	"errors"
	"fmt"

	// "io"
	"mime/multipart"
	"net/http"

	// "os"
	"time"

	"github.com/w3qxst1ck/cs2-grenades/internal/data"
	"github.com/w3qxst1ck/cs2-grenades/internal/validator"
)

func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("grenadeImage")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	v := validator.New()

	if data.VlidateImage(fileHeader, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Erorrs)
		return
	}

	fileName, err := app.saveImageToStorage(file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	grenadeID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	imageUrl := fmt.Sprintf("%s%s", app.config.storageS3.DownloadUrl, fileName)

	image := &data.Image{
		Name:      fileName,
		GrenadeID: grenadeID,
		ImageURL:  imageUrl,
	}

	err = app.models.Images.Insert(image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"image": image}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	image, err := app.models.Images.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Images.Delete(image.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// deleting from S3 storage
	// err = app.deleteImagesFromStorage([]*data.Image{image})
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "image successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) saveImageToStorage(file multipart.File) (string, error) {
	fileName := fmt.Sprintf("%d.jpg", time.Now().UnixMicro())

	// err := app.uploadImageToStorage(file, fileName)
	// if err != nil {
	// 	return "", err
	// }

	return fileName, nil
}
