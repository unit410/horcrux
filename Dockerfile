FROM golang:1.17.10-buster

WORKDIR /go/src/app
COPY . .

RUN make build

