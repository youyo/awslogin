.DEFAULT_GOAL := help
Owner := youyo
Name := awslogin
Repository := "github.com/$(Owner)/$(Name)"
GithubToken := ${GITHUB_TOKEN}
Version := $(shell git describe --tags --abbrev=0)
CommitHash := $(shell git rev-parse --verify HEAD)
BuildTime := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GoVersion := $(shell go version)

## Setup
setup:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v github.com/mitchellh/gox
	go get -u -v github.com/tcnksm/ghr
	go get -u -v github.com/jstemmer/go-junit-report

## Install dependencies
deps:
	dep ensure -v

## Vet
vet:
	go tool vet -v main.go
	go tool vet -v cmd/

## Lint
lint:
	golint -set_exit_status *.go
	golint -set_exit_status cmd

## Run tests
test:
	go test -v -cover \
		-ldflags "\
			-X \"$(Repository)/cmd/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.Version=$(Version)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		$(Repository) \
		$(Repository)/cmd/$(Name) \
		$(Repository)/cmd/$(Name)/cmd

## Execute `go run`
run:
	go run \
		-ldflags "\
			-X \"$(Repository)/cmd/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.Version=$(Version)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		cmd/$(Name)/main.go ${OPTION}

## Build
build:
	gox -osarch="darwin/amd64" \
		-ldflags="\
			-X \"$(Repository)/cmd/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.Version=$(Version)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		-output="pkg/$(Name)" \
		$(Repository)/cmd/$(Name)

## Release
release:
	for arch in "darwin_amd64"; do \
		zip -j pkg/$(Name)_$$arch.zip pkg/$(Name); \
		done
	ghr -t ${GithubToken} -u $(Owner) -r $(Name) --replace $(Version) pkg/

## Remove packages
clean:
	rm -rf pkg/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps vet lint test build release clean help
