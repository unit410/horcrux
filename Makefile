all: build

build: build/horcrux

build/horcrux: cmd/horcrux/main.go $(wildcard internal/**/*.go)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux ./cmd/horcrux

test:
	go test ./...

clean:
	rm -rf build

.PHONY: clean test
