FROM golang:1.16.7

WORKDIR /go/src/app
COPY . .

RUN make build

