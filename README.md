
## setup

```sh
# module path: a unique identifier for our project
go mod init github.com/heschmat/pixel-api-go
```
When there's a valid `go.mod` file in the root of the project, it's a module.
`go.mod` ensures reproducible builds across different environments, as it keeps the exact version of dependencies.

### httprouter

- wicked fast, thanks to its use of a radix tree fro URL matching
- sends a json 404, 405 response when matchin route cannot be found (unlike `http.ServeMux`)
- it does NOT allow conflicting routes (i.e, routes that could potentially match the same request)

## Send JSON responses

NOTE: `JSON` is just text. Hence, we can write a JSON response from a GO handler in the same way we'd write any other text responses: using `w.Write()`, `io.WriteString()` or one of the `fmt.Fprint` functions.

### JSON encoding
Go's `encoding/json` package:
For the purpose of sending JSON in an HTTP response, using `json.Marshal()` is generally the better choice.
It returns a JSON representation of the GO value in a `[]byte` slice.
```go
func Marshal(v any) ([]byte, error) {}
```

NOTE:
- A `[]byte` slice is encoded to `Base64-encoded JSON string`, not a JSON array.
- Go `time.Time` values - actually a struct behind the scenes - will be ecnoded to a formatted JSON string, and not as a JSON object.
- Any pointer values will encode as the value pointed to.

### omitzero vs. omitempty (struct tags)

Use `omitempty` if you want to omit empty slices or maps from the JSON entirely, instead of having them encode to an empty JSON array like `[]`.
Otherwise, go with `omitzero`.

### enveloping
It's valuable to think about formatting upfrong and to maintain a clear and consistend response structure across your different API endpoints.

enveloping: Including a key name (like "movie") at the top-level of the JSON.
```go
// cmd/api/helpers.go
type envelope map[string]any

// internal/data.go
type Movie struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Runtime   int       `json:"runtime,omitzero"` // omit if value is zero
	Genres    []string  `json:"genres,omitzero"`  // omit if value is empty
}

// cmd/api/movies.go
movie = data.Movie{
    ID:        id,
    Title:     "Black Swan",
    Runtime:   108,
    Genres:    []string{"drama", "thriller"},
}

movie = {"movie": movie}
```

### advanced JSON customization

So far: struct tags, adding whitespaces & enveloping the data. ✅

For more advanced JSON customization, remember this is how Go handles `JSON encoding` behind the scenes:
> When Go is encoding a particular type to JSON, it looks to see if the type satisfies the `json.Marshaler` interface. If it does, then Go will call the method, `MarshalJSON()`, to determine how to encode it, and use the `[]byte` slice that it returns as the encoded JSON value.

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
```
Hence, to customize how sth is encoded, we need to implement a `MarshalJSON()` method on it, which returns a custom JSON representation of itself in a `[]byte` slice.

For example, `time.Time` is actually a struct, but it has a `MarshalJSON()` method which outputs a string representation of itself.

In our case, to customize the `Runtime` field, a clean and simple approach is to create a custom type specifically for the field, and implement a `MarshalJSON()` method on this custom type.
The downside of having a _custom type_, when integrating our code with other packages, is that we may need to perform type conversions to change our custom type to and fro a value that other packages understand and accept.

NOTE: as **Effective Go** mentions:
> The rule about pointers vs. values for receivers is that: value methods can be invoked on pointers and values, but pointer methods can only be invoked on pointers.

## parse JSON requests

To decode a JSON form from an HTTP request body, using `json.Decoder` is generally the best choice. (compared to `json.Unmarshal()`)

```go
// 📢 When calling `.Decode()`
// we MUST pass a **non-nil pointer** as the target decode destination
// OTHERWISE, it will return `json.InvalidUnmarshalError` error at runtime.
// the struct fields MUST be exported (start with a capital letter), to be visible to `encoding/json` package.
// any JSON key-value pairs which cannot be successfully mapped to the struct fields
// - based on the struct tag names, will be silently ignored.
// `http.Server` automatically closes `r.Body`
err := json.NewDecoder(r.Body).Decode(&input)
```

Q: How can you tell the difference between a client not providing a key-value pair, and providing a key-value pair but deliberately setting it to its zero value?



## Makefile
It contains _recipes_ for automating common administra

## curl

```sh
curl localhost:4000/v1/healthcheck
# OPTIONS: what HTTP methods and features do you support for this resource?
curl -i -X OPTIONS localhost:4000/v1/healthcheck

# Method Not Allowed
curl -i -d "" localhost:4000/v1/healthcheck
curl -i -X POST localhost:4000/v1/healthcheck

curl -I localhost:4000/v1/healthcheck
curl -X HEAD localhost:4000/v1/healthcheck
curl -v -I localhost:4000/v1/healthcheck

# movies -----
curl -i localhost:4000/v1/movies/19

BODY='{"title":"Creed I","year":2015,"runtime":133,"genres":["Drama","Sport"]}'
BODY='{"title":"The Matrix","year":1999,"runtime":136,"genres":["Action","Sci-Fi"]}'
BODY='{"title":"Adore","year":2013,"runtime":111,"genres":["Drama","Romance"]}'
BODY='{"title":"Black Swan","year":2010,"runtime":108,"genres":["Drama","Thriller"]}'
BODY='{"title":"Whiplash","year":2014,"runtime":107,"genres":["Drama","Music"]}'
BODY='{"title":"Nightcrawler","year":2014,"runtime":117,"genres":["Crime","Drama","Thriller"]}'
BODY='{"title":"Her","year":2013,"runtime":126,"genres":["Drama","Romance","Sci-Fi"]}'

curl -i -d "$BODY" localhost:4000/v1/movies
```

## internal

### data
Encapsulates ALL the custom data types for our project along with the logic for interacting with our database.


## panic recovery

This is what happens if there's a `runtime panic` in our handler code:
1- Normal execution of the code in the handler will immediately stop.
2- Any deferred functions for the current goroutine will be run in reverse order, LIFO.
3- The panic will then be recovered by Go's `http.Server`, which will close the underlying HTTP connection.
4- An error message and stack trace will be output to `http.Server.ErrorLog`.
5- No other HTTP requests will be affected by the panic.

HOWEVER, it'd be better to also send a `500 Internal Server Error` response to the client, rather than just closing the HTTP connection without providing any explanation.
---
## errors
`errors.Is(err, target)` asks: is this error equal to this known value, maybe wrapped?
```go
// because io.EOF is a specific error value:
errors.Is(err, io.EOF)
```
Use it for sentinel errors like `io.EOF`, `context.Canceled`, `sql.ErrNoRows`.

`errors.As(err, &target)` asks: does this error contain a value of this type, maybe wrapped?
```go
// because json.SyntaxError is a type with useful info:
var syntaxError *json.SyntaxError
errors.As(err, &syntaxError)
```
Use it when you need fields/methods from the specific error type: `syntaxError.Offset`


## CRUD
In this section:
- How to create a **database model** that isolates all the logic for executing SQL queries against the database.

- How to implement the basic CRUD operations on a specific resource in the context of an API.

### setup the `movie` model
It will encapsulate all the code for reading and writing movie data to and from our PostgreSQL database.

We'll use `database/sql` package to execute our database queries, rather than using a 3rd-party `ORM`.


## MiSK
```sh
go install github.com/rakyll/hey@latest

# download the latest  v1.N.N release
go get github.com/julienschmidt/httprouter@v1

# download pq: a popular, reliable, and well-tested driver for PostgreSQL
go get github.com/lib/pq@v1
```

### API versioning
Real-world APIs change over time, sometimes in a backward-incompatible way.
Hence, to prevent problems for clients, it's always a good idea to implement some form of API versioning.


### fmt

A useful rule of thumb is:
- S = returns a String
- F = writes to a File-like `io.Writer`
- No prefix = writes to standard output (`os.Stdout`)

`fmt.Sprintf` formats a string and returns it
`fmt.Fprintf` formats a string and writes it to an `io.Writer`

Use `fmt.Fprintf` when you want to write directly to something that implements `io.Writer`, such as:
- an `http.ResponseWriter`
- `os.Stdout`
- a file
- a network connection
- a `bytes.Buffer`

```go
// func Fprintf(w io.Writer, format string, a ...any) (n int, err error)
// n: number of bytes written

file, _ := os.Create("output.txt")
defer file.Close()

fmt.Fprintf(file, "Hello %s\n", name)

// ------
var buf bytes.Buffer

fmt.Fprintf(&buf, "Hello %s!", "Alice")

fmt.Println(buf.String())
```

```go
fmt.Fprintf(w, "show movie with id %d", id)

// dump the contents of the input struct in an HTTP response.
fmt.Fprintf(w, "%+v\n", input)

// -----
// with %w, it wraps another error:
fmt.Errorf("read JSON: %w", err)
```

### http
```go
http.NotFound(w, r)

// -----
err = app.writeJSON(w, http.StatusOK, movie, nil)
if err != nil {
    app.logger.Error(err.Error())
    http.Error(w, "the server could not process your request", http.StatusInternalServerError)
}
```

### strconv
`strconv.Atoi`, `strconv.Quote`, ``

### httprouter
`httprouter` allows us to set our own custom error handlers. The handlers must satisfy the `http.Handler` interface.

```go
type Handler interface {
 ServeHTTP(ResponseWriter, *Request)
}
```
Hence, `http.HandlerFunc()` may be handy.
