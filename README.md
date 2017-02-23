Bytemark command-line client
============================

[![Build Status](https://travis-ci.org/BytemarkHosting/bytemark-client.svg)](https://travis-ci.org/BytemarkHosting/bytemark-client) [![Coverage Status](https://coveralls.io/repos/github/BytemarkHosting/bytemark-client/badge.svg?branch=develop)](https://coveralls.io/github/BytemarkHosting/bytemark-client?branch=develop)

Installation
------------

If you're just looking to get started you can find the most recent stable
release upon our [download page](https://github.com/BytemarkHosting/bytemark-client/releases).

If you wish to build, and install, the latest stable-release from source you can do so with:

    go get github.com/BytemarkHosting/bytemark-client/cmd/bytemark


Tracking Development
====================

If you prefer to track our in-development branch you can do that via:

    go get -d github.com/BytemarkHosting/bytemark-client/cmd/bytemark
    cd $GOPATH/src/github.com/BytemarkHosting/bytemark-client/
    git checkout develop
    cd cmd/bytemark
    go build


If you have problems to report, or contributions to make, feel free to [use the issue-tracker](https://github.com/BytemarkHosting/bytemark-client/issues)

Compatibility Guarantee
=======================

In go, semantic versioning is pretty hard. We guarantee that the following types of changes will not occur between minor versions within the lib package and all packages under its path (e.g. lib/brain)

* Publicly-exported functions and methods will not be removed, renamed, nor will their prototypes change.
* Publicly-exported struct-type fields will not be removed, renamed, nor will their types change.
* Publicly-exported variables and constants will not be removed, renamed, nor will their types change.

It's suggested that you avoid using struct embedding or interface composition with multiple types if any of those types are from bytemark-client/lib or any packages inside - bytemark-client's types are wont to have extra fields and methods added.

Breaking API change
===================

The following breaking API change to the 'lib' package occurred in version 2.0.

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

If you require the old API for longer you can use gopkg.in/BytemarkHosting/bytemark-client.v1/lib to refer to the package prior to this change.
