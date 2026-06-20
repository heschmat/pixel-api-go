package main

import (
	"encoding/json"
	"errors"
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
