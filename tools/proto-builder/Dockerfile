FROM golang:1.24-alpine

# System setup
RUN apk update && apk upgrade \
    && apk add git curl build-base autoconf automake libtool \
    && apk add protoc protobuf-dev

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

ENV PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /home/protos

#No Root
RUN adduser -D -u 1000 protos -h  /home/protos
USER protos

RUN mkdir import
RUN mkdir export
