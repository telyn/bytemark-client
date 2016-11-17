(Approximately) how to add a command
====================================

lib
----

* Add a function to call the API to a relevant file in `lib/`. The
  files are named after what resource they work on. Basically just
  whatever the last chunk of the URL is in the API call, e.g.
  `virtual_machines`, `groups`, `accounts`. Replace `_` with `-`.

* If you need to add a new type, add it to the package under lib 
  that matches the API server you're talking to.

* Add the function you just made to the interface to
  `lib/interface.go`. Now 
  `go test github.com/BytemarkHosting/bytemark-client/...` will break.

* Add a unit test for your API call to the relevant `_test.go` file
  in `lib/`.

* Add a mock version of the function to `mocks/bytemark-client.go`
  Now `go test github.com/BytemarkHosting/bytemark-client/...` will work.
  Commit if you like.

cmd/bytemark
------------

* In the init function for `cmd/bytemark/<base>.go`, append the command to
  the package-level `commands` variable, or the `adminCommands` variable if
  it is an admin-only command. To see how flags and arguments are used, take
  a look at other commands like `reimage` and `resize disk`. Read 
  `cmd/bytemark/with.go` and `cmd/bytemark/context.go` to find out more about
  how `With` works, what `Providers` are available, and what functions are
  available for reading flags out of the `Context`.

* Add a unit test for the command to `cmd/bytemark/<base>_test.go`. You're
  trying to ensure that your function turns its arguments into the
  right parameters to the API-calling function you made at the
  beginning.

Finishing up
------------

Now do testing. `go test github.com/BytemarkHosting/bytemark-client/lib` and
`go test github.com/BytemarkHosting/bytemark-client/main`. If that works then YOU'RE NOT DONE YET

Run `make` and then try out your new `bytemark`. Make sure the API
calls you're making actually work (--debug-level=5). Make sure errors
mostly don't cause panics.
