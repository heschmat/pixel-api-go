package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/heschmat/pixel-api-go/internal/data"
	"github.com/heschmat/pixel-api-go/internal/validator"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		// http.NotFound(w, r)
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
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

	// initialize a new validator
	v := validator.New()

	// if any checks fail, return a 422 Unprocessable Entity response containing the errors
	if data.ValidateMovie(v, &movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	_ = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully delete"}, nil)
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// extract the movie ID from the URL path
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// fetch the movie with the above id
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// holds the expected data from the client
	// N.B. if a field is `nil` after parsing JSON,
	// we'd know no corresponding key-value pair was provided in the JSON request body.
	var input struct {
		Title   *string  `json:"title"`
		Year    *int     `json:"year"`
		Runtime *int     `json:"runtime"`
		Genres  []string `json:"genres"` // slices already have the zero-value nil
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		app.logger.Warn("oops")
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = data.Runtime(*input.Runtime)
	}
	if input.Genres != nil {
		movie.Genres = input.Genres // no need to "dereference" a slice
	}

	// initialize a new validator
	v := validator.New()

	// if any checks fail, return a 422 Unprocessable Entity response containing the errors
	if data.ValidateMovie(v, &movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// update
	movie, err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 3)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
}
