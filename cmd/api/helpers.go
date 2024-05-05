package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	param := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(param.ByName("id"), 10, 64)
	if err != nil || id < 0 {
		return 0, errors.New("invalid id parameter")
	}

	return id, err
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// ошибка в оформлении JSON объекта
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON %d", syntaxError.Offset)

		// ошибка в оформлении JSON объекта
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// ошибка в типе объекта JSON
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// обработка ошибки пустого JSON объекта
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// обработка ошибки по лишнему полю в объекте JSON (т.к. у нас стоит ограничение dec.DisallowUnknownFields())
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			filedName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", filedName)

		// обработка ощибки по превышению входящего объекта JSON
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

			// for anything else, return the error message as-is
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contains a single JSON value")
	}
	return nil
}
