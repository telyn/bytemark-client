Setting up a workspace
======================

This is a go project so you'll need to make sure you have the following things before you can even get the code:

* a go compiler (the golang-go package in debian)
* a go workspace

Install a go compiler, then create a folder (anywhere, but ~/projects/go is a pretty common one).
Set your GOPATH to the full path to the folder you made (put it in your .profile, .bashrc, .zshrc, whatever)
Inside that folder, create src, bin and pkg folders.

Get the code
============

In the future we want to be able to `go get bytemark.co.uk/client` but that's not ready yet. For now:

	git clone git@dev.bytemark.co.uk:go-bigv-client $GOPATH/src/bytemark.co.uk/client

Install dependencies
====================

	go get github.com/tools/godep
	cd $GOPATH/src/bytemark.co.uk/client
	godep restore

Compile!
========

I made a Makefile for convenience so you can just run `make`, but all that does
is run `go build bytemark.co.uk/client/main` with the -o flag set so it comes out in
the current directory.

If you want to build for different operating systems you just need to install
the relavent go compile (e.g. golang-go-darwin-amd64) and then invoke `make` / `go build`
with the GOOS and GOARCH set to (in this example) darwin and amd64, respectively.

As a convenience a .app folder for use on OS X can be built using make Bytemark.app, which 
uses some stuff I built on my home-laptop. Should work on all intel macs though. All of the
source for that is included and should automatically be built from source if you're compiling
under OS X.
