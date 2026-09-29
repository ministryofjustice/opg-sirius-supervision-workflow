FROM golang:1.27-bookworm@sha256:648f440f42a0958804efb24df176f806f9d353b41f1c0627f666428e40310f6b

RUN go install gotest.tools/gotestsum@latest

ENV PACT_LOG_LEVEL=debug
RUN go install github.com/pact-foundation/pact-go/v2@v2.7.1
RUN pact-go -l DEBUG install
