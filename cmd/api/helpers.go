package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int, error) {
	// params := httprouter.CleanPath(r.URL.Path)

	// params is a slice of key-value pairs extracted from the URL path
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// headers: map[string][]string
// containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// js, err := json.Marshal(data)
	js, err := json.MarshalIndent(data, "", "  ") // for pretty-printing the JSON output, indenting it with two spaces for better readability.
	if err != nil {
		return err
	}
	js = append(js, '\n')

	// N.B. Go will NOT throw an error if we try to range over/read from a `nil` map.
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// decode the request body into the target destination.
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for the field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// happens if `dst` is NOT a non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			// as this should NEVER happen during normal user input; it is a programmer/configuration bug.
			// this MUST be picked up in development & tests long before deployment.
			panic(err)

		default:
			return err
		}
	}

	return nil
}
