CI_VERSION = $(shell git describe --tags --abbrev=8 --dirty --always --long)

LDFLAGS :=
LDFLAGS := "$(LDFLAGS) -X 'github.com/Hexta/envoy-tools/pkg/version.version=${CI_VERSION}'"

build:
	mkdir -p dist
	go build -ldflags=$(LDFLAGS) -o dist -v ./...
