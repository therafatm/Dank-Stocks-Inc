FROM golang:latest

EXPOSE 8080

RUN go get github.com/pilu/fresh

ADD . /go/src/test_worker
WORKDIR /go/src/test_worker

CMD fresh server.go
