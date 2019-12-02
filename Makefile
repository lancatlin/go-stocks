PKGS := ./cmd/web
GO := go
GOBUILD := $(GO) build
GORUN := $(GO) run
GOTEST := $(GO) test

build:
	$(GOBUILD) -o go-stocks $(PKGS)

run:
	$(GORUN) ./cmd/web

test:
	$(GOTEST) ./...
