ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
$(error "$$GOOS is not defined. If you are using Windows, try to re-make using 'GOOS=windows make ...' ")
endif
endif

PACKAGES    := $(shell go list ./... | grep -v '/vendor/' | grep -v '/crypto/ed25519/chainkd' | grep -v '/mining/tensority')
PACKAGES += 'github.com/anonimitycash/anonimitycash-classic/mining/tensority/go_algorithm'

BUILD_FLAGS := -ldflags "-X github.com/anonimitycash/anonimitycash-classic/version.GitCommit=`git rev-parse HEAD`"

MINER_BINARY32 := miner-$(GOOS)_386
MINER_BINARY64 := miner-$(GOOS)_amd64

ANONIMITYCASHD_BINARY32 := anonimitycashd-$(GOOS)_386
ANONIMITYCASHD_BINARY64 := anonimitycashd-$(GOOS)_amd64

ANONIMITYCASHCLI_BINARY32 := anonimitycashcli-$(GOOS)_386
ANONIMITYCASHCLI_BINARY64 := anonimitycashcli-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")

MINER_RELEASE32 := miner-$(VERSION)-$(GOOS)_386
MINER_RELEASE64 := miner-$(VERSION)-$(GOOS)_amd64

ANONIMITYCASHD_RELEASE32 := anonimitycashd-$(VERSION)-$(GOOS)_386
ANONIMITYCASHD_RELEASE64 := anonimitycashd-$(VERSION)-$(GOOS)_amd64

ANONIMITYCASHCLI_RELEASE32 := anonimitycashcli-$(VERSION)-$(GOOS)_386
ANONIMITYCASHCLI_RELEASE64 := anonimitycashcli-$(VERSION)-$(GOOS)_amd64

ANONIMITYCASH_RELEASE32 := anonimitycash-$(VERSION)-$(GOOS)_386
ANONIMITYCASH_RELEASE64 := anonimitycash-$(VERSION)-$(GOOS)_amd64

all: test target release-all install

anonimitycashd:
	@echo "Building anonimitycashd to cmd/anonimitycashd/anonimitycashd"
	@go build $(BUILD_FLAGS) -o cmd/anonimitycashd/anonimitycashd cmd/anonimitycashd/main.go

anonimitycashd-simd:
	@echo "Building SIMD version anonimitycashd to cmd/anonimitycashd/anonimitycashd"
	@cd mining/tensority/cgo_algorithm/lib/ && make
	@go build -tags="simd" $(BUILD_FLAGS) -o cmd/anonimitycashd/anonimitycashd cmd/anonimitycashd/main.go

anonimitycashcli:
	@echo "Building anonimitycashcli to cmd/anonimitycashcli/anonimitycashcli"
	@go build $(BUILD_FLAGS) -o cmd/anonimitycashcli/anonimitycashcli cmd/anonimitycashcli/main.go

install:
	@echo "Installing anonimitycashd and anonimitycashcli to $(GOPATH)/bin"
	@go install ./cmd/anonimitycashd
	@go install ./cmd/anonimitycashcli

target:
	mkdir -p $@

binary: target/$(ANONIMITYCASHD_BINARY32) target/$(ANONIMITYCASHD_BINARY64) target/$(ANONIMITYCASHCLI_BINARY32) target/$(ANONIMITYCASHCLI_BINARY64) target/$(MINER_BINARY32) target/$(MINER_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(MINER_BINARY32) $(MINER_BINARY32).exe
	cd target && cp -f $(ANONIMITYCASHD_BINARY32) $(ANONIMITYCASHD_BINARY32).exe
	cd target && cp -f $(ANONIMITYCASHCLI_BINARY32) $(ANONIMITYCASHCLI_BINARY32).exe
	cd target && md5sum $(MINER_BINARY32).exe $(ANONIMITYCASHD_BINARY32).exe $(ANONIMITYCASHCLI_BINARY32).exe >$(ANONIMITYCASH_RELEASE32).md5
	cd target && zip $(ANONIMITYCASH_RELEASE32).zip $(MINER_BINARY32).exe $(ANONIMITYCASHD_BINARY32).exe $(ANONIMITYCASHCLI_BINARY32).exe $(ANONIMITYCASH_RELEASE32).md5
	cd target && rm -f $(MINER_BINARY32) $(ANONIMITYCASHD_BINARY32) $(ANONIMITYCASHCLI_BINARY32) $(MINER_BINARY32).exe $(ANONIMITYCASHD_BINARY32).exe $(ANONIMITYCASHCLI_BINARY32).exe $(ANONIMITYCASH_RELEASE32).md5
	cd target && cp -f $(MINER_BINARY64) $(MINER_BINARY64).exe
	cd target && cp -f $(ANONIMITYCASHD_BINARY64) $(ANONIMITYCASHD_BINARY64).exe
	cd target && cp -f $(ANONIMITYCASHCLI_BINARY64) $(ANONIMITYCASHCLI_BINARY64).exe
	cd target && md5sum $(MINER_BINARY64).exe $(ANONIMITYCASHD_BINARY64).exe $(ANONIMITYCASHCLI_BINARY64).exe >$(ANONIMITYCASH_RELEASE64).md5
	cd target && zip $(ANONIMITYCASH_RELEASE64).zip $(MINER_BINARY64).exe $(ANONIMITYCASHD_BINARY64).exe $(ANONIMITYCASHCLI_BINARY64).exe $(ANONIMITYCASH_RELEASE64).md5
	cd target && rm -f $(MINER_BINARY64) $(ANONIMITYCASHD_BINARY64) $(ANONIMITYCASHCLI_BINARY64) $(MINER_BINARY64).exe $(ANONIMITYCASHD_BINARY64).exe $(ANONIMITYCASHCLI_BINARY64).exe $(ANONIMITYCASH_RELEASE64).md5
else
release: binary
	cd target && md5sum $(MINER_BINARY32) $(ANONIMITYCASHD_BINARY32) $(ANONIMITYCASHCLI_BINARY32) >$(ANONIMITYCASH_RELEASE32).md5
	cd target && tar -czf $(ANONIMITYCASH_RELEASE32).tgz $(MINER_BINARY32) $(ANONIMITYCASHD_BINARY32) $(ANONIMITYCASHCLI_BINARY32) $(ANONIMITYCASH_RELEASE32).md5
	cd target && rm -f $(MINER_BINARY32) $(ANONIMITYCASHD_BINARY32) $(ANONIMITYCASHCLI_BINARY32) $(ANONIMITYCASH_RELEASE32).md5
	cd target && md5sum $(MINER_BINARY64) $(ANONIMITYCASHD_BINARY64) $(ANONIMITYCASHCLI_BINARY64) >$(ANONIMITYCASH_RELEASE64).md5
	cd target && tar -czf $(ANONIMITYCASH_RELEASE64).tgz $(MINER_BINARY64) $(ANONIMITYCASHD_BINARY64) $(ANONIMITYCASHCLI_BINARY64) $(ANONIMITYCASH_RELEASE64).md5
	cd target && rm -f $(MINER_BINARY64) $(ANONIMITYCASHD_BINARY64) $(ANONIMITYCASHCLI_BINARY64) $(ANONIMITYCASH_RELEASE64).md5
endif

release-all: clean
	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -rf cmd/anonimitycashd/anonimitycashd
	@rm -rf cmd/anonimitycashcli/anonimitycashcli
	@rm -rf cmd/miner/miner
	@rm -rf target
	@rm -rf $(GOPATH)/bin/anonimitycashd
	@rm -rf $(GOPATH)/bin/anonimitycashcli
	@echo "Cleaning temp test data..."
	@rm -rf test/pseudo_hsm*
	@rm -rf blockchain/pseudohsm/testdata/pseudo/
	@echo "Cleaning sm2 pem files..."
	@rm -rf crypto/sm2/*.pem
	@echo "Done."

target/$(ANONIMITYCASHD_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/anonimitycashd/main.go

target/$(ANONIMITYCASHD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/anonimitycashd/main.go

target/$(ANONIMITYCASHCLI_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/anonimitycashcli/main.go

target/$(ANONIMITYCASHCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/anonimitycashcli/main.go

target/$(MINER_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/miner/main.go

target/$(MINER_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/miner/main.go

test:
	@echo "====> Running go test"
	@go test -tags "network" $(PACKAGES)

benchmark:
	@go test -bench $(PACKAGES)

functional-tests:
	@go test -timeout=5m -tags="functional" ./test 

ci: test functional-tests

.PHONY: all target release-all clean test benchmark
