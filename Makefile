GIT_VER := $(shell git describe --tags --always --dirty="-dev")

all: lint test

v:
    @echo "Version: ${GIT_VER}"

test:
	go test -v ./...

test-race:
	go test -race ./...

lint:
	gofmt -d -s .
	gofumpt -d -extra .
	go vet ./...
	staticcheck ./...
	revive -config ./linterconfig.toml ./...

fmt:
	gofmt -s -w .
	gofumpt -extra -w .
	go mod tidy