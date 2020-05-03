FROM golang:1-alpine

RUN apk add --no-cache make tini openssl redis

COPY go.mod go.sum /authentic/
WORKDIR /authentic
RUN go mod download
COPY . /authentic

RUN make compile-e2e

ENTRYPOINT [ "/sbin/tini", "--" ]