FROM golang:1

WORKDIR /go/src/app
COPY . .

RUN make build

