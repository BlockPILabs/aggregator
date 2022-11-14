VERSION := $(shell git describe --tags --always --long)

GOLDFLAGS += -X github.com/BlockPILabs/aggregator/version.Version=$(VERSION)
GOFLAGS = -ldflags "$(GOLDFLAGS)"
all: build

build:
	go build $(GOFLAGS) -o ./build/ ./cmd/aggregator
clean:
	rm -rf build/*
