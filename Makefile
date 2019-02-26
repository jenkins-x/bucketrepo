SHELL := /bin/bash
GO := GO15VENDOREXPERIMENT=1 go
NAME := bucketrepo
OS := $(shell uname)
MAIN_GO := ./internal
ROOT_PACKAGE := $(GIT_PROVIDER)/$(ORG)/$(NAME)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
BUILDFLAGS := ''
CGO_ENABLED = 0

all: fmt lint sec  build test 

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

test: 
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./... -test.v

install:
	GOBIN=${GOPATH}/bin $(GO) install -ldflags $(BUILDFLAGS) $(MAIN_GO)

fmt:
	@echo "FORMATTING"
	@FORMATTED=`$(GO) fmt ./...`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

clean:
	rm -rf bin release $(VENDOR_DIR)

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	go get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINT)
	@echo "VETTING"
	go vet ./... 
	@echo "LINTING"
	$(GOLINT) -set_exit_status ./... 

GOSEC := $(GOPATH)/bin/gosec
$(GOSEC):
	go get -u github.com/securego/gosec/cmd/gosec/...

.PHONY: sec
sec: $(GOSEC)
	@echo "SECURITY SCANNING"
	$(GOSEC) -fmt=csv ./...

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

