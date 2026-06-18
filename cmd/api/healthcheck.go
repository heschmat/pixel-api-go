package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	jq := `{"status":"ok","environment": %q,"version": %q}`

	w.Header().Set("Content-Type", "application/json")

	// w.Write([]byte(fmt.Sprintf(jq, app.config.env, version)))
	fmt.Fprintf(w, jq, app.config.env, version)
}
