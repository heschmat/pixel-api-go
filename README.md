
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
```

## internal

### data
Encapsulates ALL the custom data types for our project along with the logic for interacting with our database.

## MiSK
```sh
go install github.com/rakyll/hey@latest

# download the latest  v1.N.N release
go get github.com/julienschmidt/httprouter@v1

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