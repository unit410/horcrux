FROM golang:1.15.12

WORKDIR /go/src/app
COPY . .

RUN make build

