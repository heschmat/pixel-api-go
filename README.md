
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

## Makefile
It contains _recipes_ for automating common administra

## MiSK
```sh
go install github.com/rakyll/hey@latest

# download the latest  v1.N.N release
go get github.com/julienschmidt/httprouter@v1

```

### API versioning
Real-world APIs change over time, sometimes in a backward-incompatible way.
Hence, to prevent problems for clients, it's always a good idea to implement some form of API versioning.
