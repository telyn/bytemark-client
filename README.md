Bytemark command-line client
============================

Installation
------------

If you have a binary then it's easy - just run it!

If you have the source then it's tougher - you'll need a go workspace
(convention seems to be to have a projects/go folder in your home folder).
Stick all the source into $GOPATH/src/bytemark.co.uk/client

At the moment you'll also need to manually clone bytemark.co.uk/auth3/client -
we're looking forward to making that project `go get`-able but it's not yet.

Here's a shell script that will do it all for you assuming you want your GOPATH
to be ~/projects/go, and you want to aid in development.

	cd ~
	export GOPATH=$HOME/projects/go
	mkdir -p $GOPATH/{bin,src,pkg}
	mkdir -p $GOPATH/src/bytemark.co.uk/auth3
	git clone https://projects.bytemark.co.uk/git/auth-client $GOPATH/src/bytemark.co.uk/auth3/client
	git clone https://projects.bytemark.co.uk/git/bytemark-client $GOPATH/src/bytemark.co.uk/client

At this point you can run `go build -o bytemark .` to build it - or `make`.
Whatever tickles your particular boat.

Also run the tests with `go test ./...` - or `make test`. 

External contributions
----------------------
Please note that we don't accept issues or PRs in Github at the moment.

For now, you can find a Redmine issue tracker here:
https://projects.bytemark.co.uk/projects/bytemark-client
