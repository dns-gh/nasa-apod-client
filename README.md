## nasa-apod-client package

[![Go Report Card](https://goreportcard.com/badge/github.com/dns-gh/nasa-apod-client)](https://goreportcard.com/report/github.com/dns-gh/nasa-apod-client)

[![GoDoc](https://godoc.org/github.com/dns-gh/nasa-apod-client/nasaclient?status.png)]
(https://godoc.org/github.com/dns-gh/nasa-apod-client/nasaclient)

Nasa web client to make http requests to the APOD API: https://api.nasa.gov/api.html#Apod

## Motivation

Used in the Nasa Space Rocks Bot project: https://github.com/dns-gh/nasa-space-rocks-bot

## Installation

- It requires Go language of course. You can set it up by downloading it here: https://golang.org/dl/
- Install it here C:/Go.
- Set your GOPATH, GOROOT and PATH environment variables:

```
export GOROOT=C:/Go
export GOPATH=WORKING_DIR
export PATH=C:/Go/bin:${PATH}
```

- Download and install the package:

```
@working_dir $ go get github.com/dns-gh/nasa-apod-client/...
@working_dir $ go install github.com/dns-gh/nasa-apod-client/nasaclient
```

## Example

See the https://github.com/dns-gh/nasa-space-rocks-bot

## Tests

TODO

## LICENSE

See included LICENSE file.