# steemtorture
Simple torture client for Steem

# Usage

### Install

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


Looking at POST errors on TCP localhost:

`$GOBIN/steemtorture -c 2 -b 128 -v -m "POST"`
