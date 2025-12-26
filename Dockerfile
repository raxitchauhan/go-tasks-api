FROM golang:1.25-alpine AS build

WORKDIR /go/src/app

COPY . .