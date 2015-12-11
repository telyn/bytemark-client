package lib

import (
	"fmt"
	"net/url"
	"strings"
)

// BigVError is the basic error type which all errors return by the client library are subclassed from.
type BigVError struct {
	Method       string
	URL          *url.URL
	StatusCode   int
	RequestBody  string
	ResponseBody string
}

// BadNameError is returned when a VirtualMachineName / GroupName or AccountName is invalid.
type BadNameError struct {
	BigVError
	Type         string
	ProblemField string
	ProblemValue string
}

// NotFoundError is returned when an object was unable to be found - either because the caller doesn't have permission to see them or because they don't exist.
type NotFoundError struct {
	BigVError
}

// NotAuthorizedError is returned when an action was unable to be performed because the caller doesn't have permission.
type NotAuthorizedError struct {
	BigVError
}

// UnknownStatusCodeError is returned when an action caused BigV to return a strange status code that the client library wasn't expecting. Perhaps it's a protocol mismatch - try updating to the latest version of the library, otherwise file a bug report.
type UnknownStatusCodeError struct {
	BigVError
}

// BadRequestError is returned when a request was malformed. Report these as bugs.
type BadRequestError struct {
	BigVError
	Problems map[string][]string
}

type InternalServerError struct {
	BigVError
}

type NilAuthError struct {
	BigVError
}

func (e BigVError) Error() string {
	return fmt.Sprintf("HTTP %s %s returned %d\r\n", e.Method, e.URL.String(), e.StatusCode)
}

func (e UnknownStatusCodeError) Error() string {
	return fmt.Sprintf("An unexpected status code happened (report this as a bug!)\r\n%s", e.BigVError.Error())
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

func (e InternalServerError) Error() string {
	out := []string{"The API server returned an error"}
	if e.RequestBody != "" {
		out = append(out, fmt.Sprintf("It had this to say: %s", e.RequestBody))
	}
	return strings.Join(out, "\r\n")
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("404 Not found\r\n%s", e.BigVError.Error())
}

func (e NotAuthorizedError) Error() string {
	return fmt.Sprintf("403 Unauthorized\r\n%s", e.BigVError.Error())

}

func (e BadNameError) Error() string {
	return fmt.Sprintf("Invalid name: '%s' is a bad %s for a %s", e.ProblemValue, e.ProblemField, e.Type)
}

func (e NilAuthError) Error() string {
	return fmt.Sprintf("Authorisation wasn't set up. It's Telyn's fault.")
}
