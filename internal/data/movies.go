package data

import "time"

type Movie struct {
	ID        int       `json:"id"`
	CreatedAT time.Time `json:"-"` // don't include in JSON response
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`    // omit if value is zero
	Runtime   Runtime   `json:"runtime,omitzero"` // omit if value is zero
	Genres    []string  `json:"genres,omitzero"`  // omit if value is empty
	Version   int       `json:"version"`
}
