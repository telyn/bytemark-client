package lib

import (
	"fmt"
	"net/url"
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
}

// TooManyDiscsOnTheDancefloorError is returned when the API call would result in more than 8 discs being attached to a VM.
type TooManyDiscsOnTheDancefloorError struct {
	BigVError
}

func (e BigVError) Error() string {
	return fmt.Sprintf("HTTP %s %s returned %d\r\n", e.Method, e.URL.String(), e.StatusCode)
}

func (e UnknownStatusCodeError) Error() string {
	return fmt.Sprintf("An unexpected status code happened (report this as a bug!)\r\n%s", e.BigVError.Error())
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("Bad HTTP request\r\n%s", e.BigVError.Error())
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
