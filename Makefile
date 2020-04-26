.PHONY: build verify compile-e2e gotidy install-cobra clean docker docker-build docker-verify

DOCKER_TAG ?= latest

build:
	go build ./...

gotidy: build
	go fmt ./...
	go mod tidy

compile-e2e:
	go build -o dist/check ./internal/cmd/check
	go build -o dist/authentic ./authentic
	go build -o dist/echo ./internal/cmd/echo

verify: compile-e2e build
	check -workdir e2e -file main.lua

install-cobra:
	go get -u github.com/spf13/cobra/cobra

docker: | docker-build docker-verify

docker-build: gotidy
	docker build -t authentic:test-${DOCKER_TAG} -f Test.Dockerfile .

docker-verify:
	docker run --rm authentic:test-${DOCKER_TAG} /bin/sh /authentic/e2e/check.sh