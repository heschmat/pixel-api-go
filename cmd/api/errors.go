package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// helper func to send a JSON-formatted error message to the client
func (app *application) errResponse(w http.ResponseWriter, r *http.Request, errMsg string, statusCode int) {
	resp := envelope{"error": errMsg}

	// err := app.writeJSON(w, statusCode, resp, nil)
	// if err != nil {
	// 	app.logError(r, err)
	// 	w.WriteHeader(500)
	// }

	_ = app.writeJSON(w, statusCode, resp, nil)
}

// as our application grows, using specialist helpers like the followings
// to manage different kinds of errors
// will help ensure that our error messages remain consistent accross all our endpoints.

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "the server encountered a problem and could not process your request"
	app.errResponse(w, r, msg, http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errResponse(w, r, "the requested resource could not be found", http.StatusNotFound)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("method %s is not supported for this resource", r.Method)
	app.errResponse(w, r, msg, http.StatusMethodNotAllowed)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errResponse(w, r, err.Error(), http.StatusBadRequest)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict. please try again"
	app.errResponse(w, r, msg, http.StatusConflict)
}
