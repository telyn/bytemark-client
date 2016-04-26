.ig
bytemark.man

TODO: license header

..
.
.\" ---------------------------------------------------------------------------
.\" Title
.\" ---------------------------------------------------------------------------
.TH BYTEMARK 1 "23 March 2016" "Bytemark Client Version 0.7.0"
.SH NAME
bytemark \- command line client for managing services with Bytemark
.
.\" ---------------------------------------------------------------------------
.SH SYNOPSIS
.\" ---------------------------------------------------------------------------
.SY bytemark
.OP \-\-help
.OP \-\-force
.OP \-\-yubikey
.OP \-\-debug\-level num
.OP \-\-user username
.OP \-\-account account_name
.OP \-\-endpoint URL
.OP \-\-billing\-endpoint URL
.OP \-\-auth\-endpoint URL
.OP \-\-config\-dir dir
.OP \-\-yubikey\-otp string
.IR multi\-word\ command
.OP command\-specific\ options
.OP command\ arguments
.
.\" ---------------------------------------------------------------------------
.SH DESCRIPTION
.\" ---------------------------------------------------------------------------
.B bytemark
is the command line client used to interact with services hosted at or provided
by Bytemark
.UR https://www.bytemark.co.uk
.UE .
.
Supported operations are centred around the creation, management, and deletion
of virtual machines hosted at Bytemark. Each command issued through the client
will be translated to an operation in the API
.UR https://www.bytemark.co.uk/docs/api/
.UE ,
meaning the client is responsible for authenticating with Bytemark and sending
interpreted commands to a central server which determines and forwards the
request to the relevant destination.
.
.\" ---------------------------------------------------------------------------
.SH OPTIONS
.\" ---------------------------------------------------------------------------
Options are positional and global, applicable to any command,
and must be specified before the command when invoking the
.B bytemark
client. Other options specific to each command must be specified immediately
after the command and before any arguments relevant to that command.

.SY bytemark
.OP global\ options
.IR command
.OP command\ options
.OP command\ arguments
.
.SS Global Options
.BI \-\-help
Instead of executing the specified command issue the documentation for that
command. This will include documentation for each positional argument and
option.
.TP
.BI \-\-yubikey
Instructs the client that when it authenticates it must do so with yubikey
level authentication. If a token with this level of authentication is
available it will be used, otherwise the client will prompt for yubikey
input.
.TP
.BI \-\-debug\-level\  num
Configures the debug level for the execution of the following command. This
allows the client to show exactly what information is being sent and received
from the brain for debugging purposes.
.TP
.BI \-\-user\  username
Specifies which
.I user
the client must authenticate as, you may wish to specify
this explicitly when you are making an operation on an asset controlled by an
.I account
which has a different
.I account_name
to the
.I username
you wish to authenticate as.
.TP
.BI \-\-account\  account_name
This specifies, where not explicitly stated, which
.I account
the command executed should be in the context of. As with the
.B \-\-user
option this is normally used in cases where a
.I user
needs to reference assets held by an
.I account
with a different
.I account_name
to their
.I username
.TP
.BI \-\-endpoint\  URL
.TP
.BI \-\-billing\-endpoint\  URL
.TP
.BI \-\-auth\-endpoint\  URL
These options allow the client to talk to a different
.I endpoint
than the default option. The endpoints are the servers the client talks to in
order to authenticate, make queries, and issue commands and most users should
use the default unless otherwise instructed.
.TP
.BI \-\-config\-dir\  dir
The client reads in default settings and stores tokens in the configuration
directory. Unless otherwise specified by this option it will default to
.I $HOME/.bytemark
.TP
.BI \-\-yubikey\-otp\  string
Instead of being prompted for yubikey input this can be given when the
.B bytemark
program passing it in with this option.
.
.\" ---------------------------------------------------------------------------
.SH COMMANDS
.\" ---------------------------------------------------------------------------
Commands consist of one or more words which instruct the client what to do with
each invocation. Information explaining what each command does and what its
specific options are are available from the help text available by invoking
the help command.

.SY bytemark
.IR help

.I help
on its own prints the generic help text which will present the user with a list
of commands and topics which can then be further queried with the help command.
.YS
.SY bytemark
.IR help
.IR create
.IR server

For example prints the help text for the
.I create server
command detailing usage and options.
.YS

Commands are otherwise written directly written after the
.B bytemark
clause (and any global options). See
.B EXAMPLES.
.
.\" ---------------------------------------------------------------------------
.SH EXAMPLES
.\" ---------------------------------------------------------------------------
.
.SY bytemark
.IR create
.IR server
.IR --public-keys-file\  ~/.ssh/id_rsa.pub\ stoneboat.http

Creates the new server
.I stoneboat
in the
.I http
group and specifies also that the public key located at
.I ~/.ssh/id_rsa.pub
should be uploaded to the root and admin user on the server as part of the
post-install process.

.SY bytemark
.IR show
.IR server\  stoneboat.http

Print the hardware and networking configuration of a server named
.I stoneboat
in the
.I http
group - if the full qualifying name
.I name.group.account
or no group is specified then the client will search the default group and
account for the
invoking user.

.SY bytemark
.IR shutdown\  stoneboat.http

Send the shutdown signal to the server
.I stoneboat.http
as if the power button had been pressed.
.
.\" ---------------------------------------------------------------------------
.SH FILES
.\" ---------------------------------------------------------------------------
.TP
.I ~/.bytemark
All local configuration and temporary files are stored here.
.TP
.I ~/.bytemark/debug.log
A history of recent commands and their output on the terminal which can be used
when filing bug reports.
.TP
.I ~/.bytemark/token
The authentication token is written when one command is authenticated and used
for subsequent commands while it is still valid. It will only be valid for a
fixed amount of time.
.
.\" ---------------------------------------------------------------------------
.SH AUTHORS
.\" ---------------------------------------------------------------------------
The
.B bytemark
command line client was written and is currently maintained by Telyn Z. Roat.
The man page was written by Berin Smaldon.
.
.\" ---------------------------------------------------------------------------
.SH REPORTING BUGS
.\" ---------------------------------------------------------------------------
Right now you need to contact the Bytemark support team at
.B support@bytemark.co.uk
who will help verify bugs and file bug reports on your behalf.
