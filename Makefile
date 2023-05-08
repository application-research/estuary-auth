SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19
VERSION=$(shell git describe --always --tag --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
DOCKER_COMPOSE_FILE=docker-compose.yml
DOCKER_ORG=0utercore

.PHONY: all
all: build

.PHONY: build
build:
	go generate
	go build -tags netgo -ldflags="-s -w -X main.Commit=$(COMMIT) -X main.Version=$(VERSION)" -o auth

.PHONY: clean
clean:
	rm -f delta
	git submodule deinit --all -f

install:
	install -C -m 0755 auth /usr/local/bin
