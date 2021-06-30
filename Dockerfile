FROM golang:1.15.13

WORKDIR /go/src/app
COPY . .

RUN make build

