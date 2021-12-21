FROM golang:1.17.5

WORKDIR /go/src/app
COPY . .

RUN make build

