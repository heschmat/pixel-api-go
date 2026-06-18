
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
