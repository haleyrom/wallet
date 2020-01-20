GO_ON ?= GO111MODULE=on go
GO_OFF ?= GO111MODULE=off go
GO ?= $(GO_ON)
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell GO111MODULE=off $(GO) list ./...)
VETPACKAGES ?= $(shell GO111MODULE=off $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*")

TAGET_APP = receive_msg
TAGS = jsoniter
TAGS_PPROF = $(TAGS) pprof

LDFLAGS += -X "github.com/haleyrom/wallet/pkg/version.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "github.com/haleyrom/wallet/pkg/version.GitHash=$(shell git rev-parse HEAD)"

.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/server/server"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/server"

all: start

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -tags '$(TAGS)' -o $(SERVER_BIN) ./cmd/server

server:
	@$(GO) run -ldflags '$(LDFLAGS)' -tags '$(TAGS)' ./cmd/server/server.go

test:
	@go test -cover -race ./...

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/test/data ./cmd/server/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r server.$(NOW).zip "server"

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO_OFF) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -min_confidence 1.0 -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO_OFF) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO_OFF) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

.PHONY: tools
tools:
	$(GO_OFF) get golang.org/x/lint/golint
	$(GO_OFF) get github.com/client9/misspell/cmd/misspell
