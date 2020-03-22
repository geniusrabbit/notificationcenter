SHELL := /bin/bash -o pipefail
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
COMMIT_NUMBER ?= develop

TMP_BASE := .tmp
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)/bin
TMP_ETC := $(TMP)/etc
TMP_LIB := $(TMP)/lib
TMP_VERSIONS := $(TMP)/versions

unexport GOPATH
export GOPATH=$(abspath $(TMP))
export GO111MODULE := on
export GOBIN := $(abspath $(TMP_BIN))
export PATH := $(GOBIN):$(PATH)
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0

GOLINT_VERSION := d0100b6bd8b389f0385611eb39152c4d7c3a7905
GOLINT := $(TMP_VERSIONS)/golint/$(GOLINT_VERSION)
$(GOLINT):
	$(eval GOLINT_TMP := $(shell mktemp -d))
	cd $(GOLINT_TMP); go get golang.org/x/lint/golint@$(GOLINT_VERSION)
	@rm -rf $(GOLINT_TMP)
	@rm -rf $(dir $(GOLINT))
	@mkdir -p $(dir $(GOLINT))
	@touch $(GOLINT)

ERRCHECK_VERSION := v1.2.0
ERRCHECK := $(TMP_VERSIONS)/errcheck/$(ERRCHECK_VERSION)
$(ERRCHECK):
	$(eval ERRCHECK_TMP := $(shell mktemp -d))
	cd $(ERRCHECK_TMP); go get github.com/kisielk/errcheck@$(ERRCHECK_VERSION)
	@rm -rf $(ERRCHECK_TMP)
	@rm -rf $(dir $(ERRCHECK))
	@mkdir -p $(dir $(ERRCHECK))
	@touch $(ERRCHECK)

STATICCHECK_VERSION := c2f93a96b099cbbec1de36336ab049ffa620e6d7
STATICCHECK := $(TMP_VERSIONS)/staticcheck/$(STATICCHECK_VERSION)
$(STATICCHECK):
	$(eval STATICCHECK_TMP := $(shell mktemp -d))
	cd $(STATICCHECK_TMP); go get honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	@rm -rf $(STATICCHECK_TMP)
	@rm -rf $(dir $(STATICCHECK))
	@mkdir -p $(dir $(STATICCHECK))
	@touch $(STATICCHECK)

# UPDATE_LICENSE_VERSION := ce2550dad7144b81ae2f67dc5e55597643f6902b
# UPDATE_LICENSE := $(TMP_VERSIONS)/update-license/$(UPDATE_LICENSE_VERSION)
# $(UPDATE_LICENSE):
# 	$(eval UPDATE_LICENSE_TMP := $(shell mktemp -d))
# 	cd $(UPDATE_LICENSE_TMP); go get go.uber.org/tools/update-license@$(UPDATE_LICENSE_VERSION)
# 	@rm -rf $(UPDATE_LICENSE_TMP)
# 	@rm -rf $(dir $(UPDATE_LICENSE))
# 	@mkdir -p $(dir $(UPDATE_LICENSE))
# 	@touch $(UPDATE_LICENSE)

CERTSTRAP_VERSION := v1.1.1
CERTSTRAP := $(TMP_VERSIONS)/certstrap/$(CERTSTRAP_VERSION)
$(CERTSTRAP):
	$(eval CERTSTRAP_TMP := $(shell mktemp -d))
	cd $(CERTSTRAP_TMP); go get github.com/square/certstrap@$(CERTSTRAP_VERSION)
	@rm -rf $(CERTSTRAP_TMP)
	@rm -rf $(dir $(CERTSTRAP))
	@mkdir -p $(dir $(CERTSTRAP))
	@touch $(CERTSTRAP)

GOMOCK_VERSION := v1.3.1
GOMOCK := $(TMP_VERSIONS)/mockgen/$(GOMOCK_VERSION)
$(GOMOCK):
	$(eval GOMOCK_TMP := $(shell mktemp -d))
	cd $(GOMOCK_TMP); go get github.com/golang/mock/mockgen@$(GOMOCK_VERSION)
	@rm -rf $(GOMOCK_TMP)
	@rm -rf $(dir $(GOMOCK))
	@mkdir -p $(dir $(GOMOCK))
	@touch $(GOMOCK)

.PHONY: deps
deps: $(GOLINT) $(ERRCHECK) $(STATICCHECK) $(CERTSTRAP) $(GOMOCK)

.PHONY: generate-code
generate-code: ## Generate mocks for the project
	@echo "Generate mocks for the project"
	@go generate ./...

.PHONY: golint
golint: $(GOLINT)
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY:
errcheck: $(ERRCHECK)
	errcheck ./...

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	staticcheck ./...

.PHONY: lint
lint: golint vet errcheck staticcheck

.PHONY: test
test: ## Run package test
	go test -race ./...

.PHONY: godepup
godepup: ## Update current dependencies to the last version
	go get -u -v

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
