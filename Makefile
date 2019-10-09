all: build

build: build/horcrux

build/horcrux: cmd/horcrux/main.go $(wildcard internal/**/*.go)
	CGO_ENABLED=0 go build -o ./build/horcrux ./cmd/horcrux

test:
	go test ./...

clean:
	rm -rf build

.PHONY: clean test
