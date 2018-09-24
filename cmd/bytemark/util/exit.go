package util

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	auth3 "gitlab.bytemark.co.uk/auth/client"
)

// UserRequestedExit is returned when the user said 'No' to a 'yes/no' prompt.
type UserRequestedExit struct{}

func (e UserRequestedExit) Error() string {
	return "User requested exit"
}

// ExitCode is a named type for the E_* constants which are used as exit codes.
type ExitCode int

const (
	// ExitCodeUsageDisplayed is returned when some usage info / help page was displayed. Unsure whether it should == E_SUCCESS or not
	ExitCodeUsageDisplayed ExitCode = 0
	// ExitCodeSuccess is used to say everything went well
	ExitCodeSuccess = 0
	// ExitCodeTrappedInterrupt is the exit code returned when an unexpected interrupt like SIGUSR1 was trapped
	ExitCodeTrappedInterrupt = -1
	// ExitCodeClientBug is the exit code returned when bytemark-client knows it's faulty
	ExitCodeClientBug = 1
	// ExitCodeCantReadConfig is the exit code returned when we couldn't read a config variable from the disk for some reason
	ExitCodeCantReadConfig = 3
	// ExitCodeCantWriteConfig is the exit code returned when we couldn't write a config variable to the disk for some reason
	ExitCodeCantWriteConfig = 4
	// ExitCodeUserExit is the exit code returned when the user's action caused the program to terminate (usually by saying no to a prompt)
	ExitCodeUserExit = 5
	// ExitCodeWontDeletePopulated is the exit code returned when the user's requested that a group be deleted when it still had servers in
	ExitCodeWontDeletePopulated = 6
	// ExitCodeBadInput is the exit code returned when the user entered a malformed command, name, or flag.
	ExitCodeBadInput = 7
	// ExitCodeSubprocessFailed is the exit code returned when the client attempted to run a subprocess (e.g. ssh, a browser or a vpn client) but couldn't
	ExitCodeSubprocessFailed = 8

	// ExitCodeNoDefaultAccount is the exit code returned when the client couldn't determine a default account. In this situation, the user should manually specify the account to use with the --account flag or using `bytemark config set account`
	ExitCodeNoDefaultAccount = 9

	// ExitCodeUnknownError is the exit code returned when we got an error we couldn't deal with.
	ExitCodeUnknownError = 49

	// ExitCodeCantConnectAuth is the exit code returned when we were unable to establish an HTTP connection to the auth endpoint.
	ExitCodeCantConnectAuth = 50
	// ExitCodeCantConnectAPI is the exit code returned when we were unable to establish an HTTP connection to the API endpoint.
	ExitCodeCantConnectAPI = 150

	// ExitCodeAuthInternalError is the exit code returned when the auth server reported an internal error.
	ExitCodeAuthInternalError = 51
	// ExitCodeAPIInternalError is the exit code returned when the API server reported an internal error.
	ExitCodeAPIInternalError = 152

	// ExitCodeCantParseAuthResponse is the exit code returned when the auth server returned something we were unable to parse.
	ExitCodeCantParseAuthResponse = 52
	// ExitCodeCantParseAPIResponse is the exit code returned when the API server returned something we were unable to parse.
	ExitCodeCantParseAPIResponse = 152

	// ExitCodeInvalidCredentials is the exit code returned when the auth server says your credentials contain invalid characters.
	ExitCodeInvalidCredentials = 53
	// ExitCodeBadCredentials is the exit code returned when the auth server says your credentials don't match a user in its database.
	ExitCodeBadCredentials = 54

	// ExitCodeActionNotPermitted is the exit code returned when the API server says you haven't got permission to do that.
	ExitCodeActionNotPermitted = 155

	// ExitCodeNotFound is the exit code returned when the API server says you do not have permission to see the object you are trying to view, or that it does not exist.
	ExitCodeNotFound = 156

	// ExitCodeBadRequest is the exit code returned when we send a bad request to API. (E.g. names being too short or having wrong characters in)
	ExitCodeBadRequest = 157

	// ExitCodeUnknownAuthError is the exit code returned when we get an unexpected error from the auth server.
	ExitCodeUnknownAuthError = 149
	// ExitCodeUnknownAPIError is the exit code returned when we get an unexpected error from the Bytemark API.
	ExitCodeUnknownAPIError = 249
)

// HelpForExitCodes prints readable information on what the various exit codes do.
func HelpForExitCodes() ExitCode {
	log.Logf(`bytemark exit code list:

Exit code ranges:
    All these ranges are inclusive (i.e. 0-99 means numbers from 0 to 99, including 0 and 99.)

      0- 49: local problems
     50-149: problem talking to auth.
    150-249: problem talking to Bytemark.
    250-255: interrupts & signals

    Exit codes between 50 and 249 with the same tens and units have the same meaning but for a different endpoint

  0 -  49 Exit codes:

    0
	Nothing went wrong and I feel great!
    1
    	Problem with the client itself
    3
	Couldn't read a file from config directory
    4
	Couldn't write a file to config directory
    5
	The user caused the program to exit (usually by saying "no" to Yes/no prompts)	
    6
	The user requested a non-empty group be deleted
    7
	The program was called with malformed arguments
    8
	Attempting to execute a subprocess failed

 50 - 249 Exit codes:

     50 / 150
        Unable to establish a connection to auth/API endpoint
    
     51 / 151
        Auth endpoint reported an internal error
    
     52 / 152
        Unable to parse output from auth endpoint (probably implies a protocol mismatch - try updating bytemark)

     53
	Your credentials were rejected for containing invalid characters or fields.

     54
	Your credentials did not match any user on file - check you entered them correctly

    155
	Your user account doesn't have authorisation to perform that action

    156
        Something couldn't be found by the API server. This could be due to the following reasons:
            * It doesn't exist
	    * Your user account doesn't have authorisation to see it
	    * Protocol mismatch between the Bytemark endpoint and our client (i.e. client out of date).

    149 / 249

        An unknown error fell out of the auth / API library.

250 - 255 Exit codes:

    255
	Trapped an interrupt signal, so exited.
`)
	return ExitCodeUsageDisplayed
}

// ProcessError processes the given error, outputs a message, and returns the relevant ExitCode for the given error.
func ProcessError(err error, message ...string) ExitCode {
	if err == nil {
		return ExitCodeSuccess
	}

	trace := make([]byte, 4096)
	runtime.Stack(trace, false)

	log.Debug(log.LvlOutline, "ProcessError called. Dumping arguments and stacktrace", os.Args, string(trace))
	if len(message) > 0 {
		log.Error(message)
	}
	errorMessage := "Unknown error"
	exitCode := ExitCode(ExitCodeUnknownError)
	if err != nil {
		switch e := err.(type) {
		case *auth3.Error:
			// TODO(telyn): I feel like this entire chunk should be in gitlab.bytemark.co.uk/auth/client
			switch e.Err.(type) {
			case *url.Error:
				urlErr, _ := e.Err.(*url.Error)
				if urlErr.Err != nil {
					if opError, ok := urlErr.Err.(*net.OpError); ok {
						errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %v", opError.Err)
					} else {
						errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %T %v\r\nPlease file a bug report quoting this message.", urlErr.Err, urlErr.Err)
					}
				} else {
					errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %v", urlErr)
				}
				exitCode = ExitCodeCantConnectAuth
			default:
				errorMessage = fmt.Sprintf("Couldn't create auth session: %v", e.Err)
				exitCode = ExitCodeUnknownAuthError
			}
		case *url.Error:
			if e.Err != nil {
				if opError, ok := e.Err.(*net.OpError); ok {
					errorMessage = fmt.Sprintf("Couldn't connect to the Bytemark API: %v", opError.Err)
				} else {
					errorMessage = fmt.Sprintf("Couldn't connect to the Bytemark API: %T %v\r\nPlease file a bug report quoting this message.", e.Err, e.Err)
				}
			} else {
				errorMessage = fmt.Sprintf("Couldn't connect to the Bytemark API: %v", e)
			}
		case *exec.Error:
			if e.Name == "xdg-open" || e.Name == "x-www-browser" {
				errorMessage = "Unable to find a browser to start. You may wish to install xdg-open (part of the xdg-utils package on Debian systems)"
			} else if e.Name == "open" {
				errorMessage = "Unable to find a browser to start. Ensure that the 'open' tool is in your PATH (usually lives in /usr/bin)."

			} else if e.Name == "ssh" {
				errorMessage = fmt.Sprintf("Unable to find an SSH client, please check you have one installed.")
			} else {
				errorMessage = fmt.Sprintf("Unable to find %s in your PATH.", e.Name)
			}
			exitCode = ExitCodeSubprocessFailed
		case SubprocessFailedError:
			if e.Err == nil {
				return ExitCodeSuccess
			}
			errorMessage = err.Error()
			exitCode = ExitCodeSubprocessFailed
		case lib.NoDefaultAccountError:
			errorMessage = err.Error()
			exitCode = ExitCodeNoDefaultAccount
		case lib.NilAuthError:
			errorMessage = "Authorization wasn't set up in the client - please file a bug report containing the name of the command you tried to run."
			exitCode = ExitCodeClientBug
		case lib.ForbiddenError:
			errorMessage = err.Error()
			exitCode = ExitCodeActionNotPermitted
		case lib.BadRequestError:
			errorMessage = err.Error()
			exitCode = ExitCodeBadRequest
		case lib.ServiceUnavailableError:
			errorMessage = err.Error()
			exitCode = ExitCodeCantConnectAPI
		case lib.InternalServerError:
			errorMessage = err.Error()
			exitCode = ExitCodeAPIInternalError
		case lib.NotFoundError:
			errorMessage = err.Error()
			exitCode = ExitCodeNotFound
		case util.WontDeleteGroupWithVMsError:
			errorMessage = err.Error()
			exitCode = ExitCodeWontDeletePopulated
		case UserRequestedExit:
			errorMessage = ""
			exitCode = ExitCodeUserExit
		case *syscall.Errno:
			errorMessage = fmt.Sprintf("A command we tried to execute failed. The operating system gave us the error code %d", e)
			exitCode = ExitCodeUnknownError
		case lib.AmbiguousKeyError:
			exitCode = ExitCodeBadInput
			errorMessage = err.Error()
		case UsageDisplayedError:
			exitCode = ExitCodeBadInput
			errorMessage = err.Error()

		default:
			if fmt.Sprintf("%T", err) == "*errors.errorString" {
				errorMessage = err.Error()
				exitCode = ExitCodeBadInput // just going with BadInput because most errorStrings come from auth or validation functions.
			}
			msg := err.Error()
			if strings.Contains(msg, "Badly-formed parameters") {
				exitCode = ExitCodeInvalidCredentials
				errorMessage = "The supplied credentials contained invalid characters - please try again"
			} else if strings.Contains(msg, "Bad login credentials") {
				exitCode = ExitCodeBadCredentials
				errorMessage = "A user account with those credentials could not be found. Check your details and try again"

			}
		}

		if _, ok := err.(lib.APIError); ok && exitCode == ExitCodeUnknownError {
			errorMessage = fmt.Sprintf("Unknown error from API client library. %s", err.Error())
			exitCode = ExitCodeUnknownAPIError
		}
	} else {
		exitCode = 0
	}

	if exitCode == ExitCodeUnknownError {
		log.Errorf("Unknown error of type %T: %s.\r\nPlease send a bug report containing %s to support@bytemark.co.uk.\r\n", err, err, log.LogFile.Name())
	} else if len(message) == 0 { // the message (passed as argument) is shadowed by errorMessage (made in this function)
		log.Log(errorMessage)

	}
	return exitCode
}
