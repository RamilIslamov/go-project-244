.PHONY: help build run clean test cover lint fmt vet tidy tools

SHELL := /usr/bin/env bash

APP              := gendiff
BIN_DIR          := bin
BIN              := $(BIN_DIR)/$(APP)

GOOS             ?= $(shell go env GOOS)
GOARCH           ?= $(shell go env GOARCH)
CGO_ENABLED      ?= 0

GOLANGCI_VERSION ?= v1.60.3
PKG              := ./...
COVER_PROFILE    := coverage.out

.DEFAULT_GOAL := help
GOHOSTOS := $(shell go env GOOS)

ifeq ($(GOHOSTOS),windows)
  BIN := $(BIN_DIR)/$(APP).exe
  RACE_FLAG :=
  CGO_TEST := 0
else
  BIN := $(BIN_DIR)/$(APP)
  RACE_FLAG := -race
  CGO_TEST := 1
endif

help:
	@echo "targets: build run clean test cover lint fmt vet tidy tools"

# ===== core =====
build: $(BIN)

$(BIN):
	@mkdir -p "$(BIN_DIR)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		cd code && go build -ldflags="-s -w" -o "$(BIN)" ./cmd/gendiff/main.go

run: build
	./"$(BIN)"

clean:
	@rm -rf "$(BIN_DIR)" "$(COVER_PROFILE)"

# ===== quality =====
tidy:
	cd code && go mod tidy

fmt:
	cd code && go fmt $(PKG)

vet:
	cd code && go vet $(PKG)

test:
	cd code && CGO_ENABLED=$(CGO_TEST) go test $(PKG) $(RACE_FLAG) -count=1

cover:
	cd code && CGO_ENABLED=$(CGO_TEST) go test $(PKG) $(RACE_FLAG) \
		-coverprofile="$(COVER_PROFILE)" -covermode=atomic
	@go tool cover -func="$(COVER_PROFILE)" | tail -n 1

tools:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "install golangci-lint $(GOLANGCI_VERSION)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION); \
	}

lint: tools
	@golangci-lint version
	cd code && golangci-lint run ./...