#!/usr/bin/env bash

set -e

URL="http://localhost:4000/v1/movies"

movies=(
  '{"title":"Game Night","year":2018,"runtime":100,"genres":["Action","Comedy","Crime"]}'
  '{"title":"Before Midnight","year":2013,"runtime":109,"genres":["Drama","Romance"]}'
  '{"title":"The Dark Knight","year":2008,"runtime":152,"genres":["Action","Crime","Drama"]}'
  '{"title":"Nightmare Alley","year":2021,"runtime":150,"genres":["Crime","Drama","Thriller"]}'
  '{"title":"Two Night Stand","year":2014,"runtime":86,"genres":["Comedy","Romance"]}'
)

for body in "${movies[@]}"; do
  echo "Creating movie: $body"

  curl -i \
    -H "Content-Type: application/json" \
    -d "$body" \
    "$URL"

  echo
  echo "----------------------------------------"
done

#chmod +x add_movies.sh
#./add_movies.sh
