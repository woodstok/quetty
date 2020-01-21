GO             ?= go
GOOS           ?= $(word 1, $(subst /, " ", $(word 4, $(shell go version))))

MAKEFILE       := $(realpath $(lastword $(MAKEFILE_LIST)))
ROOT_DIR       := $(shell dirname $(MAKEFILE))
SOURCES        := $(wildcard *.go src/*.go src/*/*.go) $(MAKEFILE)


all: bin/quetty

bin:
	mkdir -p $@

test: $(SOURCES) ## Run the tests
	SHELL=/bin/sh GOOS= $(GO) test -v -tags "$(TAGS)" \
				github.com/woodstok/quetty \
				github.com/woodstok/quetty/src \

clean:
	$(RM) -r bin

bin/quetty: $(SOURCES) ## build the binary
	$(GO) build $(BUILD_FLAGS) -o $@

update: ## Update the files
	$(GO) get -u
	$(GO) mod tidy
.PHONY: help
help: ## Display this help section
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z/0-9_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help
.PHONY: test  clean  update
