FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125

RUN go install gotest.tools/gotestsum@latest

ENV PACT_LOG_LEVEL=debug
RUN go install github.com/pact-foundation/pact-go/v2@v2.7.1
RUN pact-go -l DEBUG install
