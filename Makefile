SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOLINT=${GOPATH}/bin/golangci-lint
GORELEASER=/usr/local/bin/goreleaser
GOIMPI=${GOPATH}/bin/impi
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for simple-config")
fmt:
	@go fmt ./...
lint:
	./lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOTEST) ./... -coverprofile coverage.md fmt

.PHONY: install-req fmt lint tidy test imports