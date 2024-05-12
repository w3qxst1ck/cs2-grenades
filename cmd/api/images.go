package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/w3qxst1ck/cs2-grenades/internal/data"
)

func saveImage(file multipart.File, imagesDir string) (string, error) {
	fileName := fmt.Sprintf("%d.jpg", time.Now().UnixMicro())

	dst, err := os.Create(fmt.Sprintf("%s%s", imagesDir, fileName))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return fileName, nil
}

func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("grenadeImage")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	fileName, err := saveImage(file, app.config.imagesDir)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	grenadeID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	imageUrl := fmt.Sprintf("localhost:%d%s%s", app.config.port, app.config.imagesUrl, fileName)

	image := &data.Image{
		Name: fileName,
		GrenadeID: grenadeID,
		ImageURL: imageUrl,
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
