# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

WORKDIR /go/src/app


# Copy the local package files to the container's workspace.
ADD . .

RUN go get -v -u github.com/gorilla/mux
RUN go get -v -u github.com/lib/pq

RUN go build app.go

EXPOSE 8000
EXPOSE 5432
EXPOSE 4445
EXPOSE 8888
ENTRYPOINT ["./app"]
