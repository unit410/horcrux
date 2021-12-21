all: build

build: build/horcrux-linux build/horcrux-darwin-amd64 build/horcrux-darwin-arm64

build/horcrux-linux: cmd/horcrux/main.go $(wildcard internal/**/*.go)
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-linux  ./cmd/horcrux
build/horcrux-darwin-amd64: cmd/horcrux/main.go $(wildcard internal/**/*.go)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-darwin-amd64 ./cmd/horcrux
build/horcrux-darwin-arm64: cmd/horcrux/main.go $(wildcard internal/**/*.go)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-darwin-arm64  ./cmd/horcrux

test:
	go test ./...

clean:
	rm -rf build

.PHONY: clean test