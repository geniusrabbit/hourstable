export GO111MODULE := on
export GOSUMDB := off
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0
export GOPRIVATE=sum.golang.org/*

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run -v ./...

.PHONY: test
test: ## Run package test
	go test -race ./...

.PHONY: tidy
tidy: ## Run mod tidy
	@echo "Run mod tidy"
	go mod tidy

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
