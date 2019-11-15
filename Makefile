# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: welkin welkin-cli all test clean

GOBIN = $(shell pwd)/build/bin
GO ?= latest

welkin:
	build/env.sh go run build/ci.go install ./cmd/welkin
	@echo "Done building."
	@echo "Run \"$(GOBIN)/welkin\" to launch your welkin app."

welkin-cli:
	build/env.sh go run build/ci.go install ./cmd/welkin-cli
	@echo "Done building."
	@echo "Run \"$(GOBIN)/welkin-cli\" to launch the welkin command line interface."

all:
	build/env.sh go run build/ci.go install

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*
