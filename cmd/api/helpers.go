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
	// limit the size of the request body to 1_048_576 bytes (1MB)
	// every subsequent read from r.Body is constrained by this limit.
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	// initialize the `json.Decoder`
	dec := json.NewDecoder(r.Body)
	// if the JSON from client includes any field that cannot be mapped to the target destination,
	// the decoder will return an error instead of simply ignoring the field (default behavior).
	dec.DisallowUnknownFields()

	// decode the request body to the destination
	err := dec.Decode(dst)

	// Triaging the Decoder error:
	// at this point in our application build, the `.Decode()` method could potentially return
	// the following types of error:
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// check if the error has the type `*json.SyntaxError`
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// these occur when the JSON value is the wrong type for the target destination, dst.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for the field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// if the JSON contains a field that cannot be mapped to the target destination:
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must NOT be larger than %d bytes", maxBytesError.Limit)

		// happens if `dst` is NOT a non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			// as this should NEVER happen during normal user input; it is a programmer/configuration bug.
			// this MUST be picked up in development & tests long before deployment.
			panic(err)

		default:
			return err
		}
	}

	// call .Decode() again, using a pointer to **an empty anonymous struct** as the destination.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON value only")
	}

	return nil
}
