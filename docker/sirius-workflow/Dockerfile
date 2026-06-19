FROM golang:1.26-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS base
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM node:lts-alpine3.24@sha256:5fa278c599dbba0c8f873d8717d50ecbb57c5ae6a53b7ab240c25135e0b65995 AS asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY package.json .
COPY yarn.lock .

RUN yarn install --frozen-lockfile --ignore-scripts
RUN yarn build

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/cosmtrek/air@v1.52.1 && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY internal internal
COPY web/template web/template
COPY main.go .
COPY main_test.go .

RUN go mod download

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-workflow

FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3 zlib

COPY --from=build-env /go/bin/opg-sirius-workflow opg-sirius-workflow
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-workflow"]
