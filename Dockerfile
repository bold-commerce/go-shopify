FROM golang:1.14-alpine

ENV CGO_ENABLED=0
RUN mkdir -p /go/src/github.com/bold-commerce/go-shopify
WORKDIR /go/src/github.com/bold-commerce/go-shopify
