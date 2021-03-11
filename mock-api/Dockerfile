FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get
RUN go install

CMD ["app"]