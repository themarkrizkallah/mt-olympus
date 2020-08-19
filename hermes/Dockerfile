FROM golang:1.14.2

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/themarkrizkallah/hermes/
WORKDIR $GOPATH/src/github.com/themarkrizkallah/hermes/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o hermes
