package lib

import (
	"fmt"
	"net/url"
)

type BigVError struct {
	Method       string
	URL          *url.URL
	StatusCode   int
	RequestBody  string
	ResponseBody string
}

type NotFoundError struct {
	BigVError
}

type NotAuthorizedError struct {
	BigVError
}

type UnknownStatusCodeError struct {
	BigVError
}
type BadRequestError struct {
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
