package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("Starting server on port %d", app.config.port)

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}