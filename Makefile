# Simple Makefile for a Go project

all: build install run

build:
	@echo "Building the project..."

	@go build

install:
	@echo "Installing the binary..."
	@go install

run:
	@echo "Running the command..."
	@Tic-Tac-Toe game

tests:
	@echo "Running the tests"
	@go test ./test

lint:
	@echo "Running the linter..."
	@golangci-lint run

.PHONY: all build install run
