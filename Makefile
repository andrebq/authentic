.PHONY: build verify compile-e2e gotidy setup-dev clean docker docker-build docker-verify package watch

DOCKER_TAG ?= latest

build:
	go build ./...

package:
	pkger -o res

gotidy: build
	go fmt ./...
	go mod tidy

watch:
	modd

compile-e2e:
	go build -o dist/check ./internal/cmd/check
	go build -o dist/authentic ./authentic
	go build -o dist/echo ./internal/cmd/echo

verify: compile-e2e build
	check -workdir e2e -file main.lua

setup-dev:
	go get -u github.com/spf13/cobra/cobra
	go get -u github.com/markbates/pkger/cmd/pkger

docker: | docker-build docker-verify

docker-build: gotidy
	docker build -t authentic:test-${DOCKER_TAG} -f Test.Dockerfile .

docker-run: docker-build
	docker run --rm -ti authentic:test-${DOCKER_TAG} /bin/sh

docker-verify:
	docker run -e FIREBASE_WEB_APIKEY --rm authentic:test-${DOCKER_TAG} /bin/sh /authentic/e2e/check.sh
