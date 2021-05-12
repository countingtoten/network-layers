.PHONY: build test vendor

build:
	go build ./...

test:
	go test ./...

vendor:
	go mod tidy
	go mod vendor
