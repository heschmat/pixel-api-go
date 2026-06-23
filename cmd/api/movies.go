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
	// 📢 This input struct is specifically designed to represent the JSON sent by the client.
	// It's just a temporary container for decoding the request body.
	// ⚠️ We decode into a smaller input struct containing only the fields clients are allowed to provide.
	var input struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime int      `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// keep the external API input from internal business models separate
	// this makes the code safer and easier to evolve.
	// ⚠️ The API contract is separate from the database model
	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: data.Runtime(input.Runtime),
		Genres:  input.Genres,
	}

	movie, err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// include a `Location` header to inform the client which URL the newly-created resource goes to.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// dump the contents of the input struct in an HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}
