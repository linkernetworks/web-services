## command
GO           = go
GO_VENDOR    = govendor

################################################

.PHONY: all
all: build test

.PHONY: pre-build
pre-build:
	$(GO_VENDOR) sync -v
	$(GO_VENDOR) remove -v +unused

.PHONY: build
build: pre-build
	$(GO) build -v ./...

.PHONY: test
test: build
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=1m ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html

.PHONY: check
check: check-govendor check-bats check-docker

.PHONY: clean
clean:
	$(RM) -rf $(BUILD_FOLDER)
	$(GO) clean -i -r -x -cache -testcache

## check #############################

.PHONY: check-govendor
check-govendor:
	$(info check govendor)
	@[ "`which $(GO_VENDOR)`" != "" ] || (echo "$(GO_VENDOR) is missing"; false) && (echo ".. OK")

.PHONY: check-docker
check-docker:
	$(info check docker)
	@[ "`which $(DOCKER)`" != "" ] || (echo "$(DOCKER) is missing"; false) && (echo ".. OK")
