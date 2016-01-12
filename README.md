
Bytemark command-line client
============================

Installation
------------

If you have a binary then it's easy - just run it!

If you have the source then it's tougher - you'll need a go workspace (convention seems to be to have a projects/go folder in your home folder).
Stick all the source into $GOPATH/src/bytemark.co.uk/client

At the moment you'll also need to manually clone bytemark.co.uk/auth3/client - we're looking forward to making that project `go get`-able but it's not yet.

Here's a shell script that will do it all for you assuming you want your GOPATH to be ~/projects/go

	cd ~
	mkdir -p projects/go/{bin,src,pkg}
	mkdir -p projects/go/src/bytemark.co.uk/auth3
	git clone git@gitlab.bytemark.co.uk:auth/client.git projects/go/src/bytemark.co.uk/auth3/client
	git clone git@dev.bytemark.co.uk:bytemark-client projects/gosrc/bytemark.co.uk/client
	export GOPATH=$HOME/projects/go

At this point you can run `go build -o bytemark .` to build it - or `make`. Whatever tickles your particular boat.
Also run the tests with `go test ./...` - or `make test`.
