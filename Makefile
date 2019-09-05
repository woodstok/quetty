GO             ?= go
GOOS           ?= $(word 1, $(subst /, " ", $(word 4, $(shell go version))))

MAKEFILE       := $(realpath $(lastword $(MAKEFILE_LIST)))
ROOT_DIR       := $(shell dirname $(MAKEFILE))
SOURCES        := $(wildcard *.go src/*.go src/*/*.go) $(MAKEFILE)


all: bin/quetty

bin:
	mkdir -p $@

test: $(SOURCES)
	SHELL=/bin/sh GOOS= $(GO) test -v -tags "$(TAGS)" \
				github.com/woodstok/quetty/src \

clean:
	$(RM) -r bin

bin/quetty: $(SOURCES)
	$(GO) build $(BUILD_FLAGS) -o $@

update:
	$(GO) get -u
	$(GO) mod tidy

.PHONY: test  clean  update
