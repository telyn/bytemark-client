Bytemark command-line client
============================

[![Build Status](https://travis-ci.org/BytemarkHosting/bytemark-client.svg)](https://travis-ci.org/BytemarkHosting/bytemark-client) [![Coverage Status](https://coveralls.io/repos/github/BytemarkHosting/bytemark-client/badge.svg?branch=develop)](https://coveralls.io/github/BytemarkHosting/bytemark-client?branch=develop)

Installation
------------

If you're just looking to get started you can find the most recent stable
release on our [download page](https://github.com/BytemarkHosting/bytemark-client/releases).

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

The following breaking API change to the 'lib' package occurred in version 3.0

* These `lib.Client` methods have been deleted
  * `ParseVirtualMachineName`
  * `ParseGroupName`
  * `ParseAccountName`

The following changeset can be summarised as "most `lib.Client` methods now take `lib`/`brain` structs as values, rather than as pointers."

* These `lib.Client` methods now take `lib.Account` instead of `*lib.Account`
  * `RegisterNewAccount`
* These `lib.Client` methods now take `lib.GroupName` instead of `*lib.GroupName`
  * `CreateGroup`
  * `DeleteGroup`
  * `GetGroup`
  * `CreateVirtualMachine`
* These `lib.Client` methods now take `lib.VirtualMachineName` instead of `*lib.VirtualMachineName`
  * `CreateDisc`
  * `DeleteDisc`
  * `GetDisc`
  * `ResizeDisc`
  * `SetDiscIopsLimit`
  * `AddIP`
  * `DeleteVirtualMachine`
  * `GetVirtualMachine`
  * `MoveVirtualMachine`
  * `ReimageVirtualMachine`
  * `ResetVirtualMachine`
  * `RestartVirtualMachine`
  * `StartVirtualMachine`
  * `StopVirtualMachine`
  * `ShutdownVirtualMachine`
  * `UndeleteVirtualMachine`
  * `SetVirtualMachineHardwareProfile`
  * `SetVirtualMachineHardwareProfileLock`
  * `SetVirtualMachineMemory`
  * `SetVirtualMachineCores`
  * `SetVirtualMachineCDROM`
* These `lib.Client` methods now take `brain.IPCreateRequest` instead of `*brain.IPCreateRequest`
  * `AddIP`
* These `lib.Client` methods now take `brain.ImageInstall` instead of `*brain.IPCreateRequest`
  * `ReimageVirtualMachine`
* These `lib.Client` methods now take `spp.CreditCard` instead of `*spp.CreditCard`
  * `CreateCreditCard`
  * `CreateCreditCardWithToken`
* These `lib.Client` methods now return `brain.VirtualMachine` instead of `*brain.VirtualMachine`
  * `CreateVirtualMachine` (not done)
  * `GetVirtualMachine` (not done)
* These `lib.Client` methods now return `brain.Group` instead of `*brain.Group`
  * GetGroup (not done)
* These `lib.Client` methods now return `brain.Disc` instead of `*brain.Disc`
  * GetDisc (not done)
* These `lib.Client` methods now return `lib.Account` instead of `*lib.Account`
  * RegisterAccount (not done)

If you require the old API for longer you can use `gopkg.in/BytemarkHosting/bytemark-client.v2/lib` to refer to the package prior to this change.
