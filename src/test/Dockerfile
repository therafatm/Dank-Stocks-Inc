FROM golang:latest

COPY common /go/src/common
COPY test /go/src/test
WORKDIR /go/src/test

RUN go get github.com/gavv/httpexpect
RUN go get -d ./...

ENTRYPOINT tail -f /dev/null