BUILD_GO ?= go
TEST_GO ?= go

ALL_PACKAGES := bigv.io/client/lib bigv.io/client/cmd bigv.io/client

.PHONY: test update-dependencies

all: go-bigv

go-bigv:
	$(BUILD_GO) build -o go-bigv bigv.io/client

install: all
	cp go-bigv /usr/bin/go-bigv

test:
	$(TEST_GO) test $(ALL_PACKAGES)

update-dependencies: 
	$(TEST_GO) get -ut $(ALL_PACKAGES)
	godep update $(ALL_PACKAGES)
