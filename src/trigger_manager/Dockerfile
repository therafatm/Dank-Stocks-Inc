FROM golang:latest

EXPOSE 8887

ENV WAITFORIT_VERSION="v2.2.0"
RUN curl -o /usr/local/bin/waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x /usr/local/bin/waitforit

COPY common /go/src/common
COPY transaction_service /go/src/transaction_service
COPY trigger_manager /go/src/trigger_manager
WORKDIR /go/src/trigger_manager

RUN go get github.com/pilu/fresh
RUN go get ../...
