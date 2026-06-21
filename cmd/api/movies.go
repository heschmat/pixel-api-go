package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/heschmat/pixel-api-go/internal/data"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		// http.NotFound(w, r)
		app.notFoundResponse(w, r)
		return
	}

	// dummy data to test the handler
	movie := data.Movie{
		ID:        id,
		CreatedAT: time.Now(),
		Title:     "Black Swan",
		Year:      2010,
		Runtime:   108,
		Genres:    []string{"drama", "thriller"},
		// Genres:  []string{},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// app.logger.Error(err.Error())
		// http.Error(w, "the server could not process your request", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// declare an anonymous struct to hold the information. It's our *target decode destination*.
	var input struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime int      `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		// app.errResponse(w, r, err.Error(), http.StatusBadRequest)
		app.badRequestResponse(w, r, err)
		return
	}

	// dump the contents of the input struct in an HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}
