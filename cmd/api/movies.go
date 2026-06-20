package main

import (
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
