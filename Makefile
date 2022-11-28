SHELL = /bin/bash
GO         ?= go
GO_VERSION ?= $(shell $(GO) version)
GO_ARCH    ?= arm64
GO_FLAGS   ?= CGO_ENABLED=0 GOOS=linux GOARCH=$(GO_ARCH)

NAME = ics2mattermost

deps:
	@echo ">> getting dependencies"
	$(GO) mod download

build: deps
	@echo ">> building binaries"
	$(GO_FLAGS) $(GO) generate -v
	$(GO_FLAGS) $(GO) build -o $(NAME)

run: build
	$(GO) run . -l DEBUG
