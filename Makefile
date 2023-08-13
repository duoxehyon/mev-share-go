GIT_VER := $(shell git describe --tags --always --dirty="-dev")

all: lint test

v:
    @echo "Version: ${GIT_VER}"

test:
	go test -v ./...

lint:
	go fmt ./...
	go vet ./...
	staticcheck ./...