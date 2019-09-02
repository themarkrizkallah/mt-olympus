FROM golang:latest

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/themarkrizkallah/grpc_server/
COPY . $GOPATH/src/github.com/themarkrizkallah/grpc_server/
WORKDIR $GOPATH/src/github.com/themarkrizkallah/grpc_server/

COPY server.crt .
COPY server.key .

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o grpc_server
