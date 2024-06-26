VERSION := $(shell git describe --tags --always)
CONFIG_URL := https://cfg.rpchub.io/agg/op-default.json

GOLDFLAGS := -X github.com/BlockPILabs/aggregator/version.Version=$(VERSION)
GO_OPLDFLAGS := -X github.com/BlockPILabs/aggregator/config.DefaultConfigUrl=$(CONFIG_URL)
GOFLAGS = -ldflags "$(GOLDFLAGS)"


all: build

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "$(GOLDFLAGS)" -o ./build/ ./cmd/aggregator
build-op:
	CGO_ENABLED=0 go build -ldflags "$(GOLDFLAGS) $(GO_OPLDFLAGS)"  -o ./build/ ./cmd/aggregator
#build-windows:
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o ./build/ ./cmd/aggregator
#build-mac:
#	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o ./build/ ./cmd/aggregator
clean:
	rm -rf build/*
