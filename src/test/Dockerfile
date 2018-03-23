FROM golang:alpine

RUN apk add --no-cache curl \
        curl-dev \
        libcurl \
        git \
        openssl

COPY common /go/src/common
COPY test /go/src/test
WORKDIR /go/src/test

RUN go get github.com/gavv/httpexpect
RUN go get -d ./...

ENTRYPOINT tail -f /dev/null