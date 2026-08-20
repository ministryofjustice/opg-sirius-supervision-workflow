FROM golang:1.26@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS base

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/app" \
    --shell "/sbin/nologin" \
    --uid 65532 \
    app

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM node:lts-alpine3.24@sha256:5fa278c599dbba0c8f873d8717d50ecbb57c5ae6a53b7ab240c25135e0b65995 AS asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY package*.json .
COPY .npmrc .

RUN apk add --no-cache git
RUN npm ci --ignore-scripts
RUN npm run build

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app
COPY --from=asset-env /app/web/static web/static

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest

RUN mkdir -p /app/tmp && chown -R app:app /app
ENV HOME=/home/app
ENV GOCACHE=/tmp/.cache/go-build
ENV GOMODCACHE=/tmp/.cache/go-mod

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

FROM base AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM scratch

WORKDIR /go/bin

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=build-env /go/bin/opg-sirius-workflow opg-sirius-workflow
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static

USER app:app

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=3 CMD [ "/go/bin/healthcheck" ]

ENTRYPOINT ["./opg-sirius-workflow"]
