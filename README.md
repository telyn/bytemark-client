Bytemark command-line client
============================

[![Build Status](https://travis-ci.org/BytemarkHosting/bytemark-client.svg)](https://travis-ci.org/BytemarkHosting/bytemark-client) [![Coverage Status](https://coveralls.io/repos/github/BytemarkHosting/bytemark-client/badge.svg?branch=develop)](https://coveralls.io/github/BytemarkHosting/bytemark-client?branch=develop) 

Installation
------------

If you have a binary then it's easy - just run it!

If you have a go workspace you can also just `go get github.com/BytemarkHosting/bytemark-client/cmd/bytemark` if you like, to get the latest stable.

And if you want to work on the develop branch, probably the easiest way is to `go get` it and then wipe it out and clone it by hand.

`cmd/bytemark` is where 'main' is, so `cd` into there to build or use the full import path, as with go get.

IME `make test` is a bit better to use than `go test ./...` 'cause `go test` outputs about everything in the vendor dir. But if the code structure changes then `go test ./...` will always test everything, so maybe that's better.

Feel free to open issues & merge requests on the github repo at http://github.com/BytemarkHosting/bytemark-client 

Breaking API change
===================

The following breaking API change to the 'lib' package will be arriving in master in October

* the packages lib/spp, lib/brain and lib/billing are being created.
* lib.CreditCard is moving to lib/spp.CreditCard
* lib.Person is moving to lib/billing.Person
* lib.billingAccount is moving to lib/billing.Account
* lib.brainAccount is moving to lib/brain.Account
* the following types are moving from lib to the same name under lib/brain:
  * Disc
  * Group
  * ImageInstall.go
  * IP
  * IPCreateRequest
  * IPs
  * IPSpec
  * NIC
  * StoragePool
  * User
  * VirtualMachine
  * VirtualMachineSpec
* lib.VirtualMachineName is being renamed to lib.ServerName
* lib.ParseVirtualMachineName is being renamed to lib.ParseServerName

If you require the old names for longer you can use gopkg.in/BytemarkHosting/bytemark-client.v1/lib to refer to the package as it currently stands. The next bytemark-client release will be 2.0.
