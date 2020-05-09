FROM golang:1-alpine AS builder

RUN apk add --no-cache make tini
COPY go.mod go.sum /authentic/
WORKDIR /authentic
RUN go mod download
COPY . /authentic

RUN make compile-e2e

FROM golang:1-alpine

COPY --from=builder /authentic/dist/authentic /usr/bin/authentic
RUN apk add --no-cache tini
EXPOSE 8080
CMD [ "authentic", "proxy", "--tls", "/etc/authentic/certificate", "--bind", "0.0.0.0:8080", "--cookieName", "_session" ]

ENTRYPOINT [ "/sbin/tini", "--" ]