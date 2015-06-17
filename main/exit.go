package main

import (
	bigv "bigv.io/client/lib"
	auth3 "bytemark.co.uk/auth3/client"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type UserRequestedExit struct{}

func (e *UserRequestedExit) Error() string {
	return "User requested exit"
}

// ExitCode is a named type for the E_* constants which are used as exit codes.
type ExitCode int

const (
	// E_TRAPPED_INTERRUPT is the exit code returned when an unexpected interrupt like SIGUSR1 was trapped
	E_TRAPPED_INTERRUPT ExitCode = -1
	// E_CANT_READ_CONFIG is the exit code returned when we couldn't read a config variable from the disk for some reason
	E_CANT_READ_CONFIG = 3
	// E_CANT_WRITE_CONFIG is the exit code returned when we couldn't write a config variable to the disk for some reason
	E_CANT_WRITE_CONFIG = 4
	// E_USER_EXIT is the exit code returned when the user's action caused the program to terminate (usually by saying no to a prompt)
	E_USER_EXIT = 5

	// E_UNKNOWN_ERROR is the exit code returned when we got an error we couldn't deal with.
	E_UNKNOWN_ERROR = 49

	// E_CANT_CONNECT_AUTH is the exit code returned when we were unable to establish an HTTP connection to the auth endpoint.
	E_CANT_CONNECT_AUTH = 50
	// E_CANT_CONNECT_BIGV is the exit code returned when we were unable to establish an HTTP connection to the BigV endpoint.
	E_CANT_CONNECT_BIGV = 150

	// E_AUTH_REPORTED_ERROR is the exit code returned when the auth server reported an internal error.
	E_AUTH_REPORTED_ERROR = 51
	// E_BIGV_REPORTED_ERROR is the exit code returned when the BigV server reported an internal error.
	E_BIGV_REPORTED_ERROR = 152

	// E_CANT_PARSE_AUTH is the exit code returned when the auth server returned something we were unable to parse.
	E_CANT_PARSE_AUTH = 52
	// E_CANT_PARSE_BIGV is the exit code returned when the BigV server returned something we were unable to parse.
	E_CANT_PARSE_BIGV = 152

	// E_CREDENTIALS_INVALID is the exit code returned when the auth server says your credentials contain invalid characters.
	E_CREDENTIALS_INVALID = 53
	// E_CREDENTIALS_WRONG is the exit code returned when the auth server says your credentials don't match a user in its database.
	E_CREDENTIALS_WRONG = 54

	// E_NOT_AUTHORIZED_BIGV is the exit code returned when the BigV server says you haven't got permission to do that.
	E_NOT_AUTHORIZED_BIGV = 155

	// E_NOT_FOUND_BIGV is the exit code returned when the BigV server says you do not have permission to see the object you are trying to view, or that it does not exist.
	E_NOT_FOUND_BIGV = 156

	// E_UNKNOWN_AUTH is the exit code returned when we get an unexpected error from the auth server.
	E_UNKNOWN_AUTH = 149
	// E_UNKNOWN_BIGV is the exit code returned when we get an unexpected error from the BigV server.
	E_UNKNOWN_BIGV = 249
)

// HelpForExitCodes prints readable information on what the various exit codes do.
func (cmds *CommandSet) HelpForExitCodes() {
	fmt.Println(`bigv exit code list:

Exit code ranges:
    All these ranges are inclusive (i.e. 0-99 means numbers from 0 to 99, including 0 and 99.)

      0- 49: local problems
     50-149: problem talking to auth.
    150-249: problem talking to BigV.
    250-255: interrupts & signals

    Exit codes between 50 and 249 with the same tens and units have the same meaning but for a different endpoint

  0 -  49 Exit codes:

    0
	Nothing went wrong and I feel great!

    3
	Couldn't read a file from config directory

    4
	Couldn't write a file to config directory


 50 - 249 Exit codes:

     50 / 150
        Unable to establish a connection to auth/BigV endpoint
    
     51 / 151
        Auth endpoint reported an internal error
    
     52 / 152
        Unable to parse output from auth endpoint (probably implies a protocol mismatch - try updating go-bigv)

     53
	Your credentials were rejected for containing invalid characters or fields.

     54
	Your credentials did not match any user on file - check you entered them correctly

    155
	Your user account doesn't have authorisation to perform that action

    156
        Something couldn't be found on BigV. This could be due to the following reasons:
            * It doesn't exist
	    * Your user account doesn't have authorisation to see it
	    * Protocol mismatch between the BigV endpoint and go-bigv.

    149 / 249

        An unknown error fell out of the auth / bigv library.

250 - 255 Exit codes:

    255
	Trapped an interrupt signal, so exited.
`)
}

func exit(err error, message ...string) {
	if len(message) > 0 {
		fmt.Println(strings.Join(message, "\r\n"))
	} else if err == nil {
		os.Exit(0)
	}
	errorMessage := "Unknown error"
	exitCode := E_UNKNOWN_ERROR
	if err != nil {
		switch err.(type) {
		case *auth3.Error:
			// TODO(telyn): I feel like this entire chunk should be in bytemark.co.uk/auth3/client
			authErr, _ := err.(*auth3.Error)
			switch authErr.Err.(type) {
			case *url.Error:
				urlErr, _ := authErr.Err.(*url.Error)
				if urlErr.Error != nil {
					if opError, ok := urlErr.Err.(*net.OpError); ok {
						errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %v", opError.Err)
					} else {
						errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %T %v", urlErr.Err, urlErr.Err)
					}
				} else {
					errorMessage = fmt.Sprintf("Couldn't connect to the auth server: %v", urlErr)
				}
				exitCode = E_CANT_CONNECT_AUTH
			default:
				errorMessage = fmt.Sprintf("Couldn't create auth session - internal error of type %T: %v", authErr.Err, authErr.Err)
				exitCode = E_UNKNOWN_AUTH
			}
		case bigv.NotAuthorizedError:
			errorMessage = err.Error()
			exitCode = E_NOT_AUTHORIZED_BIGV

		case bigv.NotFoundError:
			errorMessage = err.Error()
			exitCode = E_NOT_FOUND_BIGV
		case *UserRequestedExit:
			errorMessage = ""
			exitCode = E_USER_EXIT

		default:
			e := err.Error()
			if strings.Contains(e, "Badly-formed parameters") {
				exitCode = E_CREDENTIALS_INVALID
				errorMessage = "The supplied credentials contained invalid characters - please try again"
			} else if strings.Contains(e, "Bad login credentials") {
				exitCode = E_CREDENTIALS_WRONG
				errorMessage = "A user account with those credentials could not be found. Check your details and try again"

			}
		}

		if _, ok := err.(bigv.BigVError); ok && exitCode == E_UNKNOWN_ERROR {
			errorMessage = fmt.Sprintf("Unknown error from BigV client library.%s", err.Error())
			exitCode = E_UNKNOWN_BIGV
		}
	} else {
		exitCode = 0
	}

	if exitCode == E_UNKNOWN_ERROR {
		fmt.Printf("Unknown error of type %T: %s. I'm going to panic now. Expect a lot of output.\r\n", err, err)
		panic("")
	} else if len(message) == 0 {
		fmt.Println(errorMessage)

	}
	os.Exit(exitCode)
}
