export GO111MODULE := on
export GOSUMDB := off
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0
export GOPRIVATE=sum.golang.org/*

# SUBDIRS := \
# 	kafka \
# 	nats \
# 	natstream \
# 	pg \
# 	redis \
# 	wrappers/concurrency \
# 	wrappers/once/bigcache \
# 	wrappers/once/redis

.PHONY: generate-code
generate-code: ## Generate mocks for the project
	@echo "Generate mocks for the project"
	@go generate ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: test
test: ## Run package test
	go test -race ./...
	# @$(foreach PKG,${SUBDIRS}, \
	# 	pushd ${PKG} > /dev/null && set -e && go test -race ./... && popd > /dev/null ; \
	# )

.PHONY: tidy
tidy: ## Run mod tidy
	@echo "Run mod tidy"
	go mod tidy
	# @$(foreach PKG,${SUBDIRS}, \
	# 	echo "Run mod tidy in ${PKG}" ; \
	# 	pushd ${PKG} > /dev/null && set -e && go mod tidy && popd > /dev/null ; \
	# )

.PHONY: godepup
godepup: ## Update current dependencies to the last version
	go get -u -v

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
