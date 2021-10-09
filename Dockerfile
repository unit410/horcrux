FROM golang:1.16.8

WORKDIR /go/src/app
COPY . .

RUN make build

