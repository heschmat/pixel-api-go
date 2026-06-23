package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAT time.Time `json:"-"` // don't include in JSON response
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`    // omit if value is zero
	Runtime   Runtime   `json:"runtime,omitzero"` // omit if value is zero
	Genres    []string  `json:"genres,omitzero"`  // omit if value is empty
	Version   int       `json:"version"`
}

// define a MovieModel struct type which wraps a `sql.DB` connection pool.
type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie Movie) (Movie, error) {
	q := `INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	err := m.DB.QueryRow(q, args...).Scan(&movie.ID, &movie.CreatedAT, &movie.Version)

	return movie, err
}
