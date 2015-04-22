

all: go-bigv-client

go-bigv-client:
	go build -o go-bigv-client bigv.io/client

install: 
	cp go-bigv-client /usr/local/bin/go-bigv-client

test:
	go test bigv.io/client/cmd bigv.io/client/lib bigv.io/client
