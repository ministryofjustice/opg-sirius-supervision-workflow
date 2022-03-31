FROM golang:1.18.0-alpine as base
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

RUN go get -u github.com/cosmtrek/air && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM node:17.8.0-alpine3.14 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY webpack.config.js .
COPY package.json .
COPY yarn.lock .

RUN yarn install
RUN yarn build

FROM base as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-workflow

FROM alpine:3.15.3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-workflow opg-sirius-workflow
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-workflow"]