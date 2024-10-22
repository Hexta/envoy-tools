CI_VERSION = $(shell git describe --tags --abbrev=8 --dirty --always --long)

LDFLAGS := -w -s
LDFLAGS := "$(LDFLAGS) -X 'github.com/Hexta/envoy-tools/pkg/version.version=${CI_VERSION}'"

BIN_NAME_SUFFIX :=

ifdef GOOS
		BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOOS)
endif

ifdef GOARCH
		BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOARCH)
endif

.PHONY: build
build:
	mkdir -p dist
	go build -ldflags=$(LDFLAGS) -o dist/envoy-tools$(BIN_NAME_SUFFIX) -v ./cmd/envoy-tools

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	@golangci-lint run -v ./...
