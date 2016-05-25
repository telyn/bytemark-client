Bytemark command-line client
============================

Installation
------------

If you have a binary then it's easy - just run it!

If you have a go workspace you can also just `go get github.com/BytemarkHosting/bytemark-client/cmd/bytemark` if you like, to get the latest stable.

And if you want to work on the develop branch, probably the easiest way is to `go get` it and then wipe it out and clone it by hand.

`cmd/bytemark` is where 'main' is, so `cd` into there to build or use the full import path, as with go get.

IME `make test` is a bit better to use than `go test ./...` 'cause `go test` outputs about everything in the vendor dir. But if the code structure changes then `go test ./...` will always test everything, so maybe that's better.

Feel free to open issues & merge requests on the github repo at http://github.com/BytemarkHosting/bytemark-client 
