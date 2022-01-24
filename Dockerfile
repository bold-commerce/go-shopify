FROM golang:1.16-alpine

ENV CGO_ENABLED=0
RUN mkdir -p /go/src/github.com/belong-inc/go-shopify
WORKDIR /go/src/github.com/belong-inc/go-shopify
