VERSION=0.0.11
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
GO111MODULE=on

all: mackerel-plugin-log-incr-rate

.PHONY: mackerel-plugin-log-incr-rate

mackerel-plugin-log-incr-rate: main.go parser.go
	go build $(LDFLAGS) -o mackerel-plugin-log-incr-rate .

linux: main.go parser.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-log-incr-rate .

check:
	go test -v ./...
	go test -race

