# tcp

A small set of Go utilities and examples demonstrating low-level HTTP-over-TCP parsing and serving.

This repository contains a simple HTTP server built on top of a custom request parser and response writer, plus small helper programs for listening to raw TCP and sending UDP messages.

ref: https://www.rfc-editor.org/rfc/rfc9112

## Project layout

- `cmd/httpserver` — HTTP server binary (uses the internal server to accept connections and a handler that implements a few routes).
- `cmd/tcplistener` — raw TCP listener that accepts a single connection and prints parsed request parts.
- `cmd/udpsender` — simple UDP client that reads from stdin and sends lines to `localhost:42069`.
- `internal/headers` — header parsing utilities.
- `internal/request` — request parsing from a reader (supports parsing request-line, headers and Content-Length body).
- `internal/response` — response writer helpers (status line, headers, chunked bodies, trailers).
- `internal/server` — small server wrapper that accepts TCP connections, uses the request parser and response writer, and invokes a Handler.

Default listening port: `42069` (hard-coded in the examples).

## Build & run

From the repository root you can build or run components directly with `go run` or `go build`.

Run the HTTP server:

```bash
go run ./cmd/httpserver
```

Or build the binary:

```bash
go build -o bin/httpserver ./cmd/httpserver
./bin/httpserver
```

Run the TCP listener (prints a single parsed request then exits):

```bash
go run ./cmd/tcplistener
```

Run the UDP sender (interactive; sends lines to localhost:42069):

```bash
go run ./cmd/udpsender
```

## Examples

Fetch the root page:

```bash
curl -v http://localhost:42069/
```

Proxy to httpbin (example):

```bash
curl -v http://localhost:42069/httpbin/get
```

## Notes and implementation details

- The project implements a custom, minimal HTTP/1.1 parser and writer for educational purposes. It understands request-lines, headers, and bodies with Content-Length; it also supports writing chunked responses and trailers.
- Header names are normalized to lowercase in `internal/headers` and validated for allowed characters.
- The server in `internal/server` uses the parser and writes the response buffer back to the connection after the handler returns.
- The module name is `tcpgo` per `go.mod`.

## Development

- Run `go build ./...` to compile all packages and check for compile errors.
- Tests can be run with `go test ./...`.