# steemtorture
Simple torture client for Steem

Works for tcp sockets, unix sockets, and websockets.


# Usage

### Install

`go get github.com/gorilla/websocket`
`go install steemtorture.go`

### Flags

| flag | default                 | description               |
|------|-------------------------|---------------------------|
| u    | "http://127.0.0.1:8090" | "url"                     |
| b    | 4096                    | "number of blocks to get" |
| c    | 64                      | "concurrency"             |
| m    | "POST"                  | "http method"             |
| v    | false                   | "verbose blocks"          |


### Example runs:

Stress testing UNIX Socket:

`$GOBIN/steemtorture -u "unix:/tmp/steem.sock" -m "GET" -c 64`

Stress testing single websocket:

`$GOBIN/steemtorture -u "ws://127.0.0.1:8090" -c 1`

Looking at POST errors on TCP localhost:

`$GOBIN/steemtorture -c 2 -b 128 -v -m "POST"`
