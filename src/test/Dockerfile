FROM golang:latest

COPY . /go/src/test
WORKDIR /go/src/test

RUN go get github.com/gavv/httpexpect
RUN go get ./...

ENTRYPOINT ["go", "test", "./..."]
