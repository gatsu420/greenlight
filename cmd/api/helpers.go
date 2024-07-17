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

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		// Example of each case will be given
		switch {
		// curl -d '<?xml version="1.0" encoding="UTF-8"?><note><to>Alex</to></note>' localhost:4000/v1/movies
		// curl -d '{"title": "Moana", }' localhost:4000/v1/movies
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %v)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed EOF")
		case errors.As(err, &unmarshalTypeError):
			// curl -d '{"title": 123}' localhost:4000/v1/movies
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %v", unmarshalTypeError.Field)
			}
			// curl -d '["foo", "bar"]' localhost:4000/v1/movies
			return fmt.Errorf("body contains incorrect JSON type (at character %v)", unmarshalTypeError.Offset)
		// curl -X POST localhost:4000/v1/movies
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %v bytes", maxBytes)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %v", fieldName)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must contain only single JSON value")
	}

	return nil
}
