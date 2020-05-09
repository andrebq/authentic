.PHONY: build verify compile-e2e gotidy setup-dev clean docker docker-build full-verify package watch

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

setup-dev:
	go get -u github.com/spf13/cobra/cobra
	go get -u github.com/markbates/pkger/cmd/pkger

docker: | package full-verify docker-build
	docker build -t authentic:${DOCKER_TAG} -f Dockerfile .
	docker tag authentic:${DOCKER_TAG} andrebq/authentic:${DOCKER_TAG}

publish: docker
	docker push andrebq/authentic:${DOCKER_TAG}

full-verify:
	docker build -t authentic:test-${DOCKER_TAG} -f Test.Dockerfile .
	docker run -e FIREBASE_WEB_APIKEY --rm authentic:test-${DOCKER_TAG} /bin/sh /authentic/e2e/check.sh
