FROM golang:alpine

RUN apk add --update --no-cache curl \
        curl-dev \
        libcurl \
        git \
        openssl \ 
        libxml2-dev \
        build-base

ENV WAITFORIT_VERSION="v2.2.0"
RUN curl -o /usr/local/bin/waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x /usr/local/bin/waitforit

COPY common /go/src/common
COPY logger_service /go/src/logger_service
WORKDIR /go/src/logger_service


RUN go get github.com/pilu/fresh
RUN go get github.com/lib/pq
RUN go get ../...