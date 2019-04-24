SHELL := /bin/bash
GO := GO111MODULE=on go
GO_NOMOD :=GO111MODULE=off go
NAME := bucketrepo
OS := $(shell uname)
MAIN_GO := ./internal
ROOT_PACKAGE := $(GIT_PROVIDER)/$(ORG)/$(NAME)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
BUILDFLAGS := ''
CGO_ENABLED = 0

all: build fmt lint sec test 

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

.PHONY: test
test: 
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./... -test.v

.PHONY: install
install:
	GOBIN=${GOPATH}/bin $(GO) install -ldflags $(BUILDFLAGS) $(MAIN_GO)

.PHONY: fmt
fmt:
	@echo "FORMATTING"
	@FORMATTED=`$(GO) fmt ./...`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

.PHONY: clean
clean:
	rm -rf bin release

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	$(GO_NOMOD) get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINT)
	@echo "VETTING"
	go vet ./... 
	@echo "LINTING"
	$(GOLINT) -set_exit_status ./... 

GOSEC := $(GOPATH)/bin/gosec
$(GOSEC):
	$(GO_NOMOD) get -u github.com/securego/gosec/cmd/gosec/...

.PHONY: sec
sec: $(GOSEC)
	@echo "SECURITY SCANNING"
	$(GOSEC) -fmt=csv ./...

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

docker: linux
	docker build -t jenkinsxio/bucketrepo:latest .

