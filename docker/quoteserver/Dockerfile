# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

WORKDIR /go/src/app


# Copy the local package files to the container's workspace.
ADD . .

RUN go get -v -u github.com/gorilla/mux

#RUN go build -o app.go

ENTRYPOINT ["./app.go"]


