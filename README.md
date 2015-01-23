bigv-client
===========

This software package is a command-line interface for the BigV "cloud" hosting
platform. More details about BigV and Bytemark (the authors) can be found at:

  * https://bigv.io
  * https://www.bytemark.co.uk

To sign up for an account, go to: https://go.bigv.io/

For licensing information, see the file LICENSE.txt

This tool allows you to create, administer and remove resources on BigV using
the command line. It's intended to be simple, fast, and easy to deploy, so
typically comes as a precompiled static binary. To make changes or build your
own copy from source, see the file INSTALL.txt

General usage
-------------

*You must sign up for an account with BigV because this tool is useful*

Once you have an account, run:

    $ bigv profile generate --username <username> --account <account-name>

The command will ask for you password, and once given, a directory ~/.bigv will
be created containing your profile details. The password is not stored on disc
by default, so you will be prompted for your password each time you invoke the
program.

Run:

    $ bigv account show

This will output the resources you currently have on BigV in this account. In
the event that your user has access to multiple accounts, you can run:

    $ bigv account list

To switch between accounts:

    $ bigv account select

To create a new VM:

    $ bigv vm new


Commands are organised into categories, and a comprehensive help system is
included in the program. To view the current list of commands, optionally
restricting to a particular category:

    $ bigv help
    $ bigv --help
    $ bigv help vm
    $ bigv vm --help


To view help for a particular command, including available options:

    $ bigv help vm new
    $ bigv vm new --help

Modes of operation
------------------

The client operates in one of three modes, which can be specified explicitly
by including one of --batch, --oneshot or --interactive in the invocation.
So far, examples have all been given in "oneshot" mode.

If you just run:

    $ bigv

You will get the interactive mode, which is a prompt that looks like:

    bigv>

You can enter commands in this mode, one at a time; each will be executed as
if in oneshot mode before returning you to the prompt. Go back to your shell by
typing 'exit'. The main advantage of this mode is that your credentials are
cached for the duration of the session, so you don't need to re-enter your
password for each command.

If a command is given in the invocation, for instance:

    $ bigv account show

Then oneshot mode is selected. In this mode, a single command is run, prompting
for any data needed to do so. The result of the command is output, then the
client exits. You can suppress the prompts, causing the client to treat missing
arguments as an error, by passing --no-prompt in the command line.

The --batch option is equivalent to --oneshot --no-prompt --output-format=yaml.
You can also pass an alternative format name as an argument, e.g., --batch=json
if you'd prefer that. The --output-format option can also be given in oneshot
and interactive modes, but it's less useful there.

Passwordless operation in batch mode
------------------------------------

Entering the password for each command can be irritating, so there are a few
ways to avoid it. One option is to store the file in your profile configuration
file (say, ~/.bigv/default/bmcloudrc). This can be achieved by running:

    $ bigv profile password store

I'd only recommend this if the file is stored on an encrypted volume, and even
then there is a risk that someone could compromise the machine while it's
running and acquire the password that way.

You can also pass the password to the client in an environment variable -
BIGV_PASSWORD. Hopefully, this can be used in the future to integrate the
bigv client with keychain systems.

Another option, especially suited to automated jobs, is to add a limited,
passwordless privilege to your user with an IP address restriction. To make the
privilege, you'd run:

    $ bigv privilege grant --level vm_admin --to <username> \
          --password_required=false
          --ip_restrictions=<ip/mask>

You'd then run any commands you like as that user, passing --no-password as an
argument to suppress the password prompt / failure.
