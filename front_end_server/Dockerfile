FROM golang:latest

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/themarkrizkallah/front_end_server/
COPY . $GOPATH/src/github.com/themarkrizkallah/front_end_server/
WORKDIR $GOPATH/src/github.com/themarkrizkallah/front_end_server/

COPY server.crt .

COPY go.mod .
COPY go.sum .

RUN go mod download

# Dependancies are vendored

RUN go build -o front_end_server
