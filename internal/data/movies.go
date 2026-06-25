package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (m MovieModel) Get(id int) (Movie, error) {
	// avoid making unncessary database call
	if id < 1 {
		return Movie{}, ErrRecordNotFound
	}

	q := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(q, id).Scan(
		&movie.ID,
		&movie.CreatedAT,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Movie{}, ErrRecordNotFound
		} else {
			return Movie{}, err
		}
	}

	return movie, nil
}

func (m MovieModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	q := `DELETE FROM movies WHERE id = $1`

	res, err := m.DB.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) Update(movie Movie) (Movie, error) {
	// in the WHERE clause: we look for a record with a specific ID & a specific version number
	// if no movie with the combined conditin (specified ID & version) is not found anymore
	// it shows that the record has been modified since we've fetched it (updated or even deleted)
	q := `UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := m.DB.QueryRow(q, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Movie{}, ErrEditConflict
		default:
			return Movie{}, err
		}
	}

	return movie, nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]Movie, error) {
	// q := `SELECT id, created_at, title, year, runtime, genres, version
	// FROM movies
	// WHERE (title ILIKE '%' || $1 || '%' OR $1 = '')
	// AND (genres @> $2::text[] OR $2 = '{}')
	// ORDER BY id`

	// ⚠️ Important detail: because this query goes through `fmt.Sprintf`,
	// the literal `%` signs in `ILIKE` need to become `%%`
	// Otherwise fmt.Sprintf treats % as formatting syntax.
	q := fmt.Sprintf(`SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE (title ILIKE '%%' || $1 || '%%' OR $1 = '')
	AND (genres @> $2::text[] OR $2::text[] = '{}')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAT,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
