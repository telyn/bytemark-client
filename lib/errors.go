package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type UnsupportedEndpointError Endpoint

func (e UnsupportedEndpointError) Error() string {
	return fmt.Sprintf("%d was not a valid endpoint choice", e)
}

type NoDefaultAccountError struct {
	InnerErr error
}

func (e NoDefaultAccountError) Error() string {
	return "Couldn't find a default BigV account - please set one using `bytemark config set account`, or specify one on the command line using the --account flag or server.group.account or group.acccount notation."
}

// APIError is the basic error type which most errors returned by the client library are subclassed from.
type APIError struct {
	Method       string
	URL          *url.URL
	StatusCode   int
	RequestBody  string
	ResponseBody string
}

func (e APIError) Error() string {
	return fmt.Sprintf("HTTP %s %s returned %d\r\n", e.Method, e.URL.String(), e.StatusCode)
}

// BadNameError is returned when a VirtualMachineName / GroupName or AccountName is invalid.
type BadNameError struct {
	APIError
	Type         string
	ProblemField string
	ProblemValue string
}

func (e BadNameError) Error() string {
	return fmt.Sprintf("Invalid name: '%s' is a bad %s for a %s", e.ProblemValue, e.ProblemField, e.Type)
}

// NotFoundError is returned when an object was unable to be found - either because the caller doesn't have permission to see them or because they don't exist.
type NotFoundError struct {
	APIError
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("404 Not found\r\n%s", e.APIError.Error())
}

// NotAuthorizedError is returned when an action was unable to be performed because the caller doesn't have permission.
type NotAuthorizedError struct {
	APIError
}

func (e NotAuthorizedError) Error() string {
	return fmt.Sprintf("403 Unauthorized\r\n%s", e.APIError.Error())

}

// UnknownStatusCodeError is returned when an action caused API to return a strange status code that the client library wasn't expecting. Perhaps it's a protocol mismatch - try updating to the latest version of the library, otherwise file a bug report.
type UnknownStatusCodeError struct {
	APIError
}

func (e UnknownStatusCodeError) Error() string {
	return fmt.Sprintf("An unexpected status code happened (report this as a bug!)\r\n%s", e.APIError.Error())
}

// BadRequestError is returned when a request was malformed.
type BadRequestError struct {
	APIError
	Problems map[string][]string
}

func newBadRequestError(ctx BigVError, response []byte) BadRequestError {
	problems := make(map[string][]string)
	err := json.Unmarshal(response, problems)
	return BadRequestError{
		context,
		problems}
}
func (e BadRequestError) Error() string {
	if len(e.Problems) == 0 {
		return fmt.Sprintf("The API told us our request was bad\r\n%s", e.ResponseBody)
	}
	out := make([]string, len(e.Problems))
	out = append(out, "Our request had some problems:")
	for k, probs := range e.Problems {
		out = append(out, fmt.Sprintf("%s:\r\n    %s", k, strings.Join(probs, "\r\n    ")))
	}
	return strings.Join(out, "\r\n")
}

type InternalServerError struct {
	APIError
}

func (e InternalServerError) Error() string {
	out := []string{"The API server returned an error"}
	if e.RequestBody != "" {
		out = append(out, fmt.Sprintf("It had this to say: %s", e.RequestBody))
	}
	return strings.Join(out, "\r\n")
}

// ServiceUnavialableError is returned by anything that makes an HTTP request resulting in a 503
type ServiceUnavailableError struct {
	APIError
}

func (e ServiceUnavailableError) Error() string {
	return fmt.Sprintf("Bytemark's API seems to be temporarily unavailable - give it another go in a few seconds, or check on http://status.bytemark.org to see if parts of the API are currently known to be down")
}

// NilAuthError is returned when a call attempts to add authentication headers to the request, but the Client.AuthSession is nil. This is always a bug as it's an issue with the code and not with anything external.
type NilAuthError struct {
	APIError
}

func (e NilAuthError) Error() string {
	return fmt.Sprintf("Authorisation wasn't set up. Please file a bug report!")
}

// AmbiguousKeyError is returned when a call to DeleteUserAuthorizedKey has an insufficiently unique

type AmbiguousKeyError struct {
	APIError
}

func (e AmbiguousKeyError) Error() string {
	return fmt.Sprint("The specified key was ambiguous - please specify the full key")
}

// AccountCreationDeferredError is returned when we get a particular response from bmbilling.
type AccountCreationDeferredError struct{}

func (e AccountCreationDeferredError) Error() string {
	return fmt.Sprintf("Account creation request accepted\r\n\r\nYour account requires a manual check, which shouldn't take long. We'll send an email when your account is ready.")
}
