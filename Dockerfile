FROM golang:1.17rc2

WORKDIR /go/src/app
COPY . .

RUN make build

