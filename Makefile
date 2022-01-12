# Ref: https://hub.docker.com/_/golang
build_image="golang:1.17-buster@sha256:c301ce41458847a02caea52b5c2fd2a18b34f33785ddd676d38defe286bb946b"
# Paths
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

all: build

build: build/horcrux-linux-amd64 build/horcrux-darwin-amd64 build/horcrux-darwin-arm64

build/horcrux-linux-amd64:
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-linux  ./cmd/horcrux
build/horcrux-darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-darwin-amd64 ./cmd/horcrux
build/horcrux-darwin-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -o ./build/horcrux-darwin-arm64  ./cmd/horcrux

build-release: clean lint
	docker run --rm \
		--platform "linux/amd64" \
		--volume $(mkfile_dir):/build \
		--workdir /build \
		$(build_image) \
		make build

test: lint
	go test $$(go list ./... | grep -v /vendor/)

coverage:
	go test $$(go list ./... | grep -v /vendor/) -coverprofile cover.out
	go tool cover -func cover.out

integration:
	./test/integration

lint:
	go mod verify
	go vet $$(go list ./... | grep -v /vendor/)
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	rm -rf build

.PHONY: all
