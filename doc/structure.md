Sort of roughly how the bytemark-client is structured as a piece of software
============================================================================

bytemark-client is sort of separated into two codebases: cmd/bytemark and lib

cmd/bytemark is where all the command-line-y stuff happens. Parsing the arguments, deciding what to do, opening the user's web browser, etc. all lives there.
lib is intended as a golang API for bytemark. It has some relatively nice things like GetVirtualMachine(name string) (*VirtualMachine, error).

Before I get into execution flow, here is a brief description of the purposes of various folders and files:

cmd/bytemark/util/config.go - this is where the 'configuration' is implemented. Basically this deals with the global flags and the ~/.bytemark folder
cmd/bytemark/util/* - utility stuff for the command-liney-bit. This is where we hide the signup form code, prompting,, extra flag types for codegangsta/cli, and the code to open the user's web browser.
cmd/bytemark/main.go - this is the entry point for the program.
cmd/bytemark/with.go - this is where With and the various Providers are defined. They're used to keep boilerplate down in the command definitions.
cmd/bytemark/context.go - this is where 
cmd/bytemark/* - where all the commands are implemented
lib/*_calls.go - all the operations that talk to the bytemark API are defined in these files
lib/*_type.go - type definitions and methods for those types belong in these files (virtualmachine_type.go is where VirtualMachine.TotalDiscSize is defined, for instance)
util/ - a dead-trivial logging library used by both parts of the client. (future work: remove its tendrils from lib somehow)
mocks/* - mock versions of the various types used by the client are defined here, to make testing easier (frankly the mocks could probably live in the same package as the thing they mock, but whatever, I had no idea what I was doing when I started this project :-D )
gen/* - scripts for generating things, used during building / releasing.
ports/* - chocolatey packaging (which should probably be in its own repo?) and a mac .App skeleton (superceded by remembering that homebrew exists)
vendor/ - other people's code

Execution flow
==============

All the init()s from cmd/bytemark run and populate cmd/bytemark.commands ([]codegangsta/cli.Command)
cmd/bytemark.main() reads the global args and sets up cmd/bytemark.global.Config (cmd/bytemark/util.Config) and cmd/bytemark.global.Client (lib.Client)
cmd/bytemark.main() sets up cmd/bytemark.global.App (a codegangsta/cli.App) using the commands array and calls Run on it.
codegansta/cli does the necessary argument-parsing and runs the relevant command.
The command then sets up authentication if necessary and calls the necessary lib.Client methods to accomplish its goals.
