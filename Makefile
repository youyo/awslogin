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
	go get -u -v github.com/Songmu/ghch/cmd/ghch
	go get -u -v github.com/motemen/gobump/cmd/gobump

## Install dependencies
deps:
	dep ensure -v

## Vet
vet:
	go vet -v $(shell go list ./...)

## Lint
lint:
	golint -set_exit_status $(shell go list ./...) && echo 'No problem.'

## Run tests
test:
	go test -v -cover \
		-ldflags "\
			-X \"$(Repository)/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" $(shell go list ./...)

## Execute `go run`
run:
	go run \
		-ldflags "\
			-X \"$(Repository)/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		$(Name)/main.go ${OPTION}

## Build
build:
	gox -osarch="darwin/amd64" \
		-ldflags="\
			-X \"$(Repository)/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		-output="pkg/$(Name)" \
		$(Repository)/$(Name)

## Release
release:
	for arch in "darwin_amd64"; do \
		zip -j pkg/$(Name)_$$arch.zip pkg/$(Name); \
		done
	ghr -t ${GithubToken} -u $(Owner) -r $(Name) --replace $(Version) pkg/

## Bump up
bump-up:
	cp CHANGELOG.md{,_old}
	$(eval ver := $(shell gobump patch -v -w | jq -r '.Version'))
	ghch --format=markdown --next-version=$(ver) > CHANGELOG.md
	cat CHANGELOG.md_old >> CHANGELOG.md
	rm -f CHANGELOG.md_old
	git add .
	git commit -m "bump up"
	git tag $(ver)

## Remove packages
clean:
	rm -rf pkg/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: help
.SILENT:
