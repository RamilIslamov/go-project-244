.PHONY: build run clean

SHELL := /usr/bin/env bash

BIN := bin/gendiff

build: $(BIN)

$(BIN):
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $(BIN) ./cmd/gendiff/main.go

run: build
	./$(BIN)

clean:
	rm -rf bin