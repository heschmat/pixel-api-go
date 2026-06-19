package main

import (
	"net/http"
)

// func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// 	jq := `{"status":"ok","environment": %q,"version": %q}`

// 	w.Header().Set("Content-Type", "application/json")

// 	// w.Write([]byte(fmt.Sprintf(jq, app.config.env, version)))
// 	fmt.Fprintf(w, jq, app.config.env, version)
// }

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	// if there's an error during the marshaling process,
	// log the error and send a 500 Internal Server Error response to the client.
	if err != nil {
		app.logger.Error("failed to marshal JSON", "error", err)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
}
