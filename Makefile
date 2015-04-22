ALL_PACKAGES := bigv.io/client/lib bigv.io/client/cmd bigv.io/client

.PHONY: test update-dependencies

all: go-bigv-client

go-bigv-client:
	go build -o go-bigv-client bigv.io/client

install: all
	cp go-bigv-client /usr/bin/go-bigv-client

test:
	go test $(ALL_PACKAGES)

update-dependencies: 
	go get -ut $(ALL_PACKAGES)
	godep update $(ALL_PACKAGES)
