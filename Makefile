

all:
	go build -o go-bigv-client bigv.io/client/main	

test:
	go test bigv.io/client/cmd bigv.io/client/lib bigv.io/client/main
