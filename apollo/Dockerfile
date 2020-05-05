FROM golang:1.13

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/themarkrizkallah/apollo/
#COPY . $GOPATH/src/github.com/themarkrizkallah/apollo/
#COPY ../proto/ $GOPATH/src/github.com/themarkrizkallah/apollo/proto/
WORKDIR $GOPATH/src/github.com/themarkrizkallah/apollo/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o apollo
