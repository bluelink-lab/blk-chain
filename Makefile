#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags))
COMMIT := $(shell git log -1 --format='%H')

BUILDDIR ?= $(CURDIR)/build
INVARIANT_CHECK_INTERVAL ?= $(INVARIANT_CHECK_INTERVAL:-0)
export PROJECT_HOME=$(shell git rev-parse --show-toplevel)
export GO_PKG_PATH=$(shell go env GOPATH)/go/pkg
export GO111MODULE = on

# process build tags

LEDGER_ENABLED ?= true
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
	ifeq ($(OS),Windows_NT)
		GCCEXE = $(shell where gcc.exe 2> NUL)
		ifeq ($(GCCEXE),)
			$(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
		else
			build_tags += ledger
		endif
	else
		UNAME_S = $(shell uname -s)
		ifeq ($(UNAME_S),OpenBSD)
			$(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
		else
			GCC = $(shell command -v gcc 2> /dev/null)
			ifeq ($(GCC),)
				$(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
			else
				build_tags += ledger
			endif
		endif
	endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=blt \
			-X github.com/cosmos/cosmos-sdk/version.ServerName=blkd \
			-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
			-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
			-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

# go 1.23+ needs a workaround to link memsize (see https://github.com/fjl/memsize).
# NOTE: this is a terribly ugly and unstable way of comparing version numbers,
# but that's what you get when you do anything nontrivial in a Makefile.
ifeq ($(firstword $(sort go1.23 $(shell go env GOVERSION))), go1.23)
	ldflags += -checklinkname=0
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

# BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)' -race
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

#### Command List ####

all: lint install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/blkd

install-with-race-detector: go.sum
		go install -race $(BUILD_FLAGS) ./cmd/blkd

install-price-feeder: go.sum
		go install $(BUILD_FLAGS) ./oracle/price-feeder

loadtest: go.sum
		go build $(BUILD_FLAGS) -o ./build/loadtest ./loadtest/

price-feeder: go.sum
		go build $(BUILD_FLAGS) -o ./build/price-feeder ./oracle/price-feeder

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		@go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

build:
	go build $(BUILD_FLAGS) -o ./build/blkd ./cmd/blkd

build-price-feeder:
	go build $(BUILD_FLAGS) -o ./build/price-feeder ./oracle/price-feeder

clean:
	rm -rf ./build

build-loadtest:
	go build -o build/loadtest ./loadtest/


###############################################################################
###                       Local testing using docker container              ###
###############################################################################
# To start a 4-node cluster from scratch:
# make clean && make docker-cluster-start
# To stop the 4-node cluster:
# make docker-cluster-stop
# If you have already built the binary, you can skip the build:
# make docker-cluster-start-skipbuild
###############################################################################


# Build linux binary on other platforms
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc make build
.PHONY: build-linux

build-price-feeder-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc make build-price-feeder
.PHONY: build-price-feeder-linux

# Build docker image
build-docker-node:
	@cd docker && docker build --tag blk-chain/localnode localnode --platform linux/x86_64
.PHONY: build-docker-node

build-rpc-node:
	@cd docker && docker build --tag blk-chain/rpcnode rpcnode --platform linux/x86_64
.PHONY: build-rpc-node

# Run a single node docker container
run-local-node: kill-blk-node build-docker-node
	@rm -rf $(PROJECT_HOME)/build/generated
	docker run --rm \
	--name blk-node \
	--user="$(shell id -u):$(shell id -g)" \
	-v $(PROJECT_HOME):/bluelink-lab/blk-chain:Z \
	-v $(GO_PKG_PATH)/mod:/root/go/pkg/mod:Z \
	-v $(shell go env GOCACHE):/root/.cache/go-build:Z \
	--platform linux/x86_64 \
	-p 26656-26658:26656-26658 \
	-p 9090-9091:9090-9091 \
	-p 8545-8546:8545-8546 \
	-p 7171:7171 \
	blk-chain/localnode
.PHONY: run-local-node

# Run a single rpc state sync node docker container
run-rpc-node: build-rpc-node
	docker run --rm \
	--name she-rpc-node \
	--network docker_localnet \
	--user="$(shell id -u):$(shell id -g)" \
	-v $(PROJECT_HOME):/bluelink-lab/blk-chain:Z \
	-v $(GO_PKG_PATH)/mod:/root/go/pkg/mod:Z \
	-v $(shell go env GOCACHE):/root/.cache/go-build:Z \
	-p 26656-26658:26656-26658 \
	--platform linux/x86_64 \
	blk-chain/rpcnode
.PHONY: run-rpc-node

run-rpc-node-skipbuild: build-rpc-node
	docker run --rm \
	--name she-rpc-node \
	--network docker_localnet \
	--user="$(shell id -u):$(shell id -g)" \
	-v $(PROJECT_HOME):/bluelink-lab/blk-chain:Z \
	-v $(GO_PKG_PATH)/mod:/root/go/pkg/mod:Z \
	-v $(shell go env GOCACHE):/root/.cache/go-build:Z \
	-p 26656-26658:26656-26658 \
	--platform linux/x86_64 \
	--env SKIP_BUILD=true \
	blk-chain/rpcnode
.PHONY: run-rpc-node

kill-blk-node:
	docker ps --filter name=blk-node --filter status=running -aq | xargs docker kill 2> /dev/null || true

kill-rpc-node:
	docker ps --filter name=she-rpc-node --filter status=running -aq | xargs docker kill 2> /dev/null || true

# Run a 4-node docker containers
docker-cluster-start: docker-cluster-stop build-docker-node
	@rm -rf $(PROJECT_HOME)/build/generated
	@mkdir -p $(shell go env GOPATH)/pkg/mod
	@mkdir -p $(shell go env GOCACHE)
	@cd docker && USERID=$(shell id -u) GROUPID=$(shell id -g) GOCACHE=$(shell go env GOCACHE) NUM_ACCOUNTS=10 INVARIANT_CHECK_INTERVAL=${INVARIANT_CHECK_INTERVAL} UPGRADE_VERSION_LIST=${UPGRADE_VERSION_LIST} docker compose up

.PHONY: localnet-start

# Use this to skip the blkd build process
docker-cluster-start-skipbuild: docker-cluster-stop build-docker-node
	@rm -rf $(PROJECT_HOME)/build/generated
	@cd docker && USERID=$(shell id -u) GROUPID=$(shell id -g) GOCACHE=$(shell go env GOCACHE) NUM_ACCOUNTS=10 SKIP_BUILD=true docker compose up
.PHONY: localnet-start

# Stop 4-node docker containers
docker-cluster-stop:
	@cd docker && USERID=$(shell id -u) GROUPID=$(shell id -g) GOCACHE=$(shell go env GOCACHE) docker compose down
.PHONY: localnet-stop


# Implements test splitting and running. This is pulled directly from
# the github action workflows for better local reproducibility.

GO_TEST_FILES != find $(CURDIR) -name "*_test.go"

# default to four splits by default
NUM_SPLIT ?= 4

$(BUILDDIR):
	mkdir -p $@

# The format statement filters out all packages that don't have tests.
# Note we need to check for both in-package tests (.TestGoFiles) and
# out-of-package tests (.XTestGoFiles).
$(BUILDDIR)/packages.txt:$(GO_TEST_FILES) $(BUILDDIR)
	go list -f "{{ if (or .TestGoFiles .XTestGoFiles) }}{{ .ImportPath }}{{ end }}" ./... | sort > $@

split-test-packages:$(BUILDDIR)/packages.txt
	split -d -n l/$(NUM_SPLIT) $< $<.
test-group-%:split-test-packages
	cat $(BUILDDIR)/packages.txt.$* | xargs go test -parallel 4 -mod=readonly -timeout=10m -race -coverprofile=$*.profile.out -covermode=atomic
