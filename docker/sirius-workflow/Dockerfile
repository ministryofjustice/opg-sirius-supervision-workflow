FROM node:lts-alpine as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY webpack.config.js .
COPY package.json .
COPY yarn.lock .

RUN yarn && yarn build

FROM golang:1.14 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-workflow

FROM alpine:3.10

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-workflow /go/bin/opg-sirius-workflow
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["/go/bin/opg-sirius-workflow"]