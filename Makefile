VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/kibelasync.revision=$(CURRENT_REVISION)"
u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps: deps
	sh -c '\
	tmpdir=$$(mktemp -d); \
	cd $$tmpdir; \
	go get ${u} \
	  golang.org/x/lint/golint            \
	  github.com/Songmu/godzil/cmd/godzil \
	  github.com/tcnksm/ghr; \
	rm -rf $$tmpdir'

.PHONY: test
test: deps
	go test

.PHONY: lint
lint: devel-deps
	golint -set_exit_status ./...

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/kibelasync

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/kibelasync

.PHONY: bump
bump: devel-deps
	godzil release

CREDITS: go.sum devel-deps
	godzil credits -w

DIST_DIR ?= dist/v$(VERSION)
.PHONY: crossbuild
crossbuild: CREDITS
	env CGO_ENABLED=0 godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
      -os=linux,darwin -d=$(DIST_DIR) ./cmd/*

.PHONY: upload
upload:
	ghr -body="$$(godzil changelog --latest -F markdown)" v$(VERSION) $(DIST_DIR)
