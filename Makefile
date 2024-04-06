.PHONY: npmi init lint-install remove test help info invoke-initiator invoke-summarizer invoke-notificator logs-initiator logs-summarizer logs-processor

PROJECT_NAME=$(shell basename "$(PWD)")
LINT_VERSION = v1.54.2

npmi:
	npm install

init: npmi
	make -C infraestructure init STACK=$(STACK)
	make init-cmd

init-cmd: npmi
	go generate ./...
	make -C cmd init STACK=$(STACK)

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

remove: npmi
	make -C cmd remove STACK=$(STACK)
	make -C infraestructure remove STACK=$(STACK)

remove-cmd: npmi
	make -C cmd remove STACK=$(STACK)

refresh: npmi
	make -C infraestructure refresh STACK=$(STACK)

test:
	go test ./... --race -coverprofile=coverage.txt -covermode=atomic

invoke-initiator:
	make -C cmd/serverless/account_summarizers invoke-initiator STACK=$(STACK)

invoke-summarizer:
	make -C cmd/serverless/account_summarizers invoke-summarizer STACK=$(STACK)

invoke-notificator:
	make -C cmd/serverless/account_summarizers invoke-notificator STACK=$(STACK)

info:
	make -C cmd/serverless/account_summarizers info STACK=$(STACK)

logs-initiator:
	make -C cmd/serverless/account_summarizers logs-initiator

logs-summarizer:
	make -C cmd/serverless/account_summarizers logs-summarizer

logs-notificator:
	make -C cmd/serverless/account_summarizers logs-notificator

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  npmi           				Install npm dependencies"
	@echo "  init STACK=local				Deploy the infrastructure, if STACK=local, deploy locally otherwise deploy in the cloud"
	@echo "  remove STACK=local				Remove the infrastructure, if STACK=local, remove locally otherwise remove in the cloud"
	@echo "  help           				Show this help message"
	@echo "  test           				Run tests"
	@echo "  info STACK=local				Show information about the stack, if STACK=local, show information about the local stack otherwise show information about the cloud stack"
	@echo "  invoke-initiator STACK=local			Invoke the initiator function, if STACK=local, invoke the local function otherwise invoke the cloud function"
	@echo "  invoke-summarizer STACK=local			Invoke the summarizer function, if STACK=local, invoke the local function otherwise invoke the cloud function"
	@echo "  invoke-notificator STACK=local			Invoke the notificator function, if STACK=local, invoke the local function otherwise invoke the cloud function"
	@echo "  logs-initiator				Show the logs of the initiator function"
	@echo "  logs-summarizer				Show the logs of the summarizer function"
	@echo "  logs-notificator				Show the logs of the notificator function"
