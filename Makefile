ALL_PACKAGES := bigv.io/client/lib bigv.io/client/cmd bigv.io/client

.PHONY: test update-dependencies

all: go-bigv

go-bigv:
	go build -o go-bigv bigv.io/client

install: all
	cp go-bigv /usr/bin/go-bigv

test:
	go test $(ALL_PACKAGES)

update-dependencies: 
	go get -ut $(ALL_PACKAGES)
	godep update $(ALL_PACKAGES)
