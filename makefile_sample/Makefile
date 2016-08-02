SOURCES=$(shell find . -name "*.go" | grep -v vendor/)
PACKAGES=$(shell go list ./... | grep -v vendor/)
FGT := $(shell command -v fgt)
GOLINT := $(shell command -v golint)
ERRCHECK := $(shell command -v errcheck)
GO_CARPET := $(shell command -v go-carpet)

default: linters test

# lint tools

linters: fmt imports lint vet errcheck
.PHONY: linters

linters-ci: linters-ci-get fmt-ci imports-ci lint-ci vet-ci errcheck-ci
.PHONY: linters-ci

linters-ci-get:
ifndef FGT
	go get -u github.com/GeertJohan/fgt
endif
ifndef GOLINT
	go get -u github.com/golang/lint/golint
endif
ifndef ERRCHECK
	go get -u github.com/kisielk/errcheck
endif
.PHONY: linters-ci-get

fmt:
	gofmt -s -w $(SOURCES)
.PHONY: fmt

fmt-ci:
	fgt gofmt -l $(SOURCES)
.PHONY: fmt-ci

imports:
	goimports -w $(SOURCES)
.PHONY: imports

imports-ci:
	fgt goimports -d -e $(SOURCES)
.PHONY: imports-ci

lint:
	go list ./... | grep -v vendor/ | grep -v Generated | xargs -L1 golint
.PHONY: lint

lint-ci:
	go list ./... | grep -v vendor/ | grep -v Generated | xargs -L1 fgt golint
.PHONY: lint-ci

vet:
	go tool vet -composites=false $(SOURCES)
.PHONY: vet

vet-ci:
	fgt go tool vet -composites=false $(SOURCES)
.PHONY: vet-ci

errcheck:
	errcheck -ignore Close $(PACKAGES)
.PHONY: errcheck

errcheck-ci:
	$(FGT) errcheck -ignore Close $(PACKAGES)
.PHONY: errcheck-ci

# testing

test:
	go test -race -v $(PACKAGES)
.PHONY: test

test-ci: linters-ci
	GORACE="halt_on_error=1" go test -race -v $(PACKAGES)
.PHONY: test-ci

test-color:
	GORACE="halt_on_error=1" go test -race -v $(PACKAGES) | \
		sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | \
		sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | \
		sed ''/RUN/s//$$(printf "\033[34mRUN\033[0m")/''
.PHONY: test-color

# code coverage

cov:
ifndef GO_CARPET
	go get -u github.com/msoap/go-carpet
endif
	go-carpet -256colors | less -R
.PHONY: cov

cov-rm:
	rm go-carpet-coverage*
.PHONY: cov-rm

# godep

deps-get:
	go get -v -t ./...
.PHONY: deps-get

deps-lock:
	rm -rf Godeps vendor
	godep save -t $(PACKAGES)
.PHONY: deps-lock

deps-restore:
	godep restore
.PHONY: deps-restore

