Bytemark command-line client
============================

Installation
------------

If you have a binary then it's easy - just run it!

If you have the source then it's tougher - you'll need a go workspace
(convention seems to be to have a projects/go folder in your home folder).
Stick all the source into $GOPATH/src/github.com/BytemarkHosting/bytemark-client

At the moment you'll also need to manually clone github.com/BytemarkHosting/auth-client -
we're looking forward to making that project `go get`-able but it's not yet.

Here's a shell script that will do it all for you assuming you want your GOPATH
to be ~/projects/go, and you want to aid in development.

	cd ~
	export GOPATH=$HOME/projects/go
	mkdir -p $GOPATH/{bin,src,pkg}
	mkdir -p $GOPATH/src/bytemark.co.uk/auth3
	git clone https://projects.bytemark.co.uk/git/auth-client $GOPATH/src/github.com/BytemarkHosting/auth-client
	git clone https://projects.bytemark.co.uk/git/bytemark-client $GOPATH/src/github.com/BytemarkHosting/bytemark-client
	cd $GOPATH/src/github.com/BytemarkHosting/bytemark-client
	make gensrc

At this point you can run `go build github.com/BytemarkHosting/bytemark-client/cmd/bytemark` to build it - or `make`.
To install it, `go install github.com/BytemarkHosting/bytemark-client/cmd/bytemark` to install it to your go folder or `make install` to install it and the manpage in /usr/bin
Whatever tickles your particular boat.

Also run the tests with `go test ./...` - or `make test`. 

External contributions
----------------------
Please note that we don't accept issues or PRs in Github at the moment.

For now, you can find a Redmine issue tracker here:
https://projects.bytemark.co.uk/projects/bytemark-client
