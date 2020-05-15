PREFIX ?= "/usr/bin/"

all: build

build:
	mkdir -p bin
	go build -o bin/jen ./cmd/jen

test: build
	go test -coverprofile=coverage.bin -race ./...
	go tool cover -func=coverage.bin

coverage: test
	go tool cover -html=coverage.bin

install: build
	install -m 755 ./bin/jen $(PREFIX)/jen

.PHONY = build test coverage install
