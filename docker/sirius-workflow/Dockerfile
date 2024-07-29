FROM golang:1.22 as base
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM node:20.16.0-alpine3.20 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY package.json .
COPY yarn.lock .

RUN yarn install
RUN yarn build

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/cosmtrek/air@v1.52.1 && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-workflow

FROM alpine:3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3

COPY --from=build-env /go/bin/opg-sirius-workflow opg-sirius-workflow
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-workflow"]
