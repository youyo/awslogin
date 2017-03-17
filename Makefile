Name := awslogin
Version := $(shell git describe --tags --abbrev=0)
OWNER := youyo
.DEFAULT_GOAL := help

## Setup
setup:
	go get github.com/kardianos/govendor
	go get github.com/Songmu/make2help/cmd/make2help
	go get github.com/mitchellh/gox

## Install dependencies
deps: setup
	govendor sync

## Vet
vet: setup
	govendor vet +local

## Lint
lint: setup
	go get github.com/golang/lint/golint
	govendor vet +local
	for pkg in $$(govendor list -p -no-status +local); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

## Run tests
test: deps
	govendor test +local -cover

## Build
build: deps
	gox -osarch="darwin/amd64 linux/amd64" -ldflags="-X main.Version=$(Version) -X main.Name=$(Name)" -output="pkg/$(Name)_{{.OS}}_{{.Arch}}"

## Build
build-local: deps
	go build -ldflags "-X main.Version=$(Version) -X main.Name=$(Name)" -o __$(Name)

## Install
install: deps
	go install -ldflags "-X main.Version=$(Version) -X main.Name=$(Name)"

## Release
release: build
	for arch in "darwin_amd64" "linux_amd64"; do \
		zip -j pkg/$(Name)_$$arch.zip pkg/$(Name)_$$arch; \
		done
	ghr -t ${GITHUB_TOKEN} -u $(OWNER) -r $(Name) --replace $(Version) pkg/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps vet lint test build build-local install release help
