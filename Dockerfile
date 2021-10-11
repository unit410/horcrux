FROM golang:1.17.2

WORKDIR /go/src/app
COPY . .

RUN make build

