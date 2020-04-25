.PHONY: build verify compile-e2e tidy install-cobra

build:
	go build ./...

tidy: build
	go fmt ./...
	go mod tidy

compile-e2e:
	go build ./internal/cmd/check
	go build ./authentic

verify: compile-e2e build
	check -workdir e2e -file main.lua

install-cobra:
	go get -u github.com/spf13/cobra/cobra