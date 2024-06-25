VERSION := $(shell git describe --tags --always)
CONFIG_URL := https://cfg.rpchub.io/agg/op-default.json

GOLDFLAGS += -X github.com/BlockPILabs/aggregator/version.Version=$(VERSION)
GOLDFLAGS += -X github.com/BlockPILabs/aggregator/config.DefaultConfigUrl=$(CONFIG_URL)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

all: build

.PHONY: build
build:
	go build $(GOFLAGS) -o ./build/ ./cmd/aggregator
clean:
	rm -rf build/*
