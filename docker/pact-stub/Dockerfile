FROM golang:1.16

WORKDIR /app

COPY pact_stub.go .

RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pact_stub

ENTRYPOINT /go/bin/pact_stub
