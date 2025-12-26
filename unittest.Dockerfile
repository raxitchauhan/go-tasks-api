FROM golang:1.25-alpine AS golang

RUN apk add --no-cache \
	build-base

RUN go env -w CGO_ENABLED=1

WORKDIR $GOPATH/src/app

COPY . .