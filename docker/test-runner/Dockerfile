FROM golang:1.27-bookworm@sha256:5cf287a799e6b94384bad13d16b14904c531f51ba65792237e122ce42b392f61

RUN go install gotest.tools/gotestsum@latest

ENV PACT_LOG_LEVEL=debug
RUN go install github.com/pact-foundation/pact-go/v2@v2.7.1
RUN pact-go -l DEBUG install
