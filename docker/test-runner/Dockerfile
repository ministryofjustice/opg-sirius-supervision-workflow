FROM golang:1.26-bookworm@sha256:b305420a68d0f229d91eb3b3ed9e519fcf2cf5461da4bef997bf927e8c0bfd2b

RUN go install gotest.tools/gotestsum@latest

ENV PACT_LOG_LEVEL=debug
RUN go install github.com/pact-foundation/pact-go/v2@v2.7.1
RUN pact-go -l DEBUG install
