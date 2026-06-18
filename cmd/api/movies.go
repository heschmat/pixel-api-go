package main

import (
	"fmt"
	"net/http"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		// http.Error(w, "invalid id parameter", http.StatusBadRequest)
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show movie with id %d", id)
}
