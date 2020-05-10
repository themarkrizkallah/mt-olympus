FROM golang:1.14.2

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/themarkrizkallah/matcher/
#COPY . $GOPATH/src/github.com/themarkrizkallah/matcher/
WORKDIR $GOPATH/src/github.com/themarkrizkallah/matcher/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o matcher
