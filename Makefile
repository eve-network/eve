#!/usr/bin/make -f


VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
ifeq ($(VERSION),)
	VERSION := v0.0.0
endif

PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8

export GO111MODULE = on

# process build tags

LEDGER_ENABLED ?= true
build_tags = netgo
build_tags += pebbledb
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

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=eved \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=eve \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb \
		  -X github.com/tendermint/tm-db.ForceSync=1 \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"
	

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
BUILDDIR ?= $(CURDIR)/build
#### Command List ####

all: install

install: go.sum
		go install $(BUILD_FLAGS) ./eved

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		@go mod verify

build:
	go build $(BUILD_FLAGS) -o ./bin/eved ./eved

# https://github.com/cosmos/ibc-go/blob/main/Makefile#L377
###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=v0.7
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=$(PROJECT_NAME)-proto-gen-$(protoVer)
containerProtoGenAny=$(PROJECT_NAME)-proto-gen-any-$(protoVer)
containerProtoGenSwagger=$(PROJECT_NAME)-proto-gen-swagger-$(protoVer)
containerProtoFmt=$(PROJECT_NAME)-proto-fmt-$(protoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen.sh; fi

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	@echo "Generating Protobuf Any"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenAny}$$"; then docker start -a $(containerProtoGenAny); else docker run --name $(containerProtoGenAny) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen-any.sh; fi

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protoc-swagger-gen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -not -path "./third_party/*" -name "*.proto" -exec clang-format -i {} \; ; fi


proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################
PACKAGES_UNIT=$(shell go list ./... | grep -E -v 'tests/simulator|e2e')
PACKAGES_E2E=$(shell go list -tags e2e ./... | grep '/e2e')
TEST_PACKAGES=./...

test: test-unit
test-unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_UNIT)

# test-e2e-ci runs a full e2e test suite
# does not do any validation about the state of the Docker environment
# As a result, avoid using this locally.
test-e2e-ci:
	@VERSION=$(VERSION) EVE_E2E_SKIP_STATE_SYNC=True EVE_E2E_SKIP_UPGRADE=True EVE_E2E_DEBUG_LOG=True go test -tags e2e -mod=readonly -timeout=25m -v $(PACKAGES_E2E) -tags e2e

# test-e2e-debug runs a full e2e test suite but does
# not attempt to delete Docker resources at the end.
test-e2e-debug: e2e-setup
	@VERSION=$(VERSION) EVE_E2E_SKIP_STATE_SYNC=True EVE_E2E_SKIP_UPGRADE=True EVE_E2E_SKIP_CLEANUP=True go test -tags e2e -mod=readonly -timeout=25m -v $(PACKAGES_E2E) -count=1 -tags e2e

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES_UNIT)

build-e2e-script:
	mkdir -p $(BUILDDIR)
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ ./tests/e2e/initialization/$(E2E_SCRIPT_NAME)

docker-build-debug:
	@DOCKER_BUILDKIT=1 docker build -t eve:${COMMIT} --build-arg BASE_IMG_TAG=debug -f tests/e2e/initialization/Dockerfile .
	@DOCKER_BUILDKIT=1 docker tag eve:${COMMIT} eve:debug

docker-build-e2e-init-chain:
	@DOCKER_BUILDKIT=1 docker build -t eve-e2e-init-chain:debug --build-arg E2E_SCRIPT_NAME=chain -f tests/e2e/initialization/init.Dockerfile .

docker-build-e2e-init-node:
	@DOCKER_BUILDKIT=1 docker build -t eve-e2e-init-node:debug --build-arg E2E_SCRIPT_NAME=node -f tests/e2e/initialization/init.Dockerfile .

.PHONY: test-mutation

###############################################################################
###                                Localnet                                 ###
###############################################################################

# Build image for a local testnet
localnet-build:
	docker build -f Dockerfile -t eve-node .

# Start a 4-node testnet locally
localnet-start: localnet-clean
	@if ! [ -f build/node0/eved/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/eve:Z eve-node -c "eved testnet --v 4 -o eve --chain-id eve-1 --keyring-backend=test --starting-ip-address 192.168.11.2"; fi
	docker-compose up -d
	bash scripts/add-keys.sh

# Clean testnet
localnet-clean:
	docker-compose down
	sudo rm -rf build

# Stop testnet
localnet-stop:
	docker-compose down


# Reset testnet
localnet-unsafe-reset:
	docker-compose down
ifeq ($(OS),Windows_NT)
	@docker run --rm -v $(CURDIR)\build\node0\eved:/eve\Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)\build\node1\eved:/eve\Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)\build\node2\eved:/eve\Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)\build\node3\eved:/eve\Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
else
	@docker run --rm -v $(CURDIR)/build/node0/eved:/eve:Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)/build/node1/eved:/eve:Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)/build/node2/eved:/eve:Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
	@docker run --rm -v $(CURDIR)/build/node3/eved:/eve:Z eve/node "./eved tendermint unsafe-reset-all --home=/eve"
endif

# Clean testnet
localnet-show-logstream:
	docker-compose logs --tail=1000 -f

.PHONY: localnet-build localnet-start localnet-stop
