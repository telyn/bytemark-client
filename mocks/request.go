package mocks

import (
	"io"
	"net/url"
	"reflect"
	"testing"
)

type Request struct {
	StatusCode     int
	ResponseBody   []byte
	ResponseObject interface{}
	Err            error
	T              *testing.T
}

func (r *Request) GetURL() url.URL {
	return url.URL{
		Scheme: "HTTP",
		Host:   "fake-host",
		Path:   "fake-path",
	}
}

func (r *Request) AllowInsecure() {

}

func (r *Request) fillOut(out interface{}) {
	if out == nil || r.ResponseObject == nil {
		return
	}
	resVal := reflect.ValueOf(r.ResponseObject)
	if resVal.Kind() == reflect.Ptr {
		resVal = resVal.Elem()
	}
	outVal := reflect.ValueOf(out)
	if resVal.Type().AssignableTo(outVal.Type()) {
		r.T.Fatalf("ResponseBody %s was not assignable to out %s", resVal.Type(), outVal.Type())
	}
	outVal.Elem().Set(resVal)
}

func (r *Request) MarshalAndRun(in interface{}, out interface{}) (statusCode int, responseBody []byte, err error) {
	r.fillOut(out)
	return r.StatusCode, r.ResponseBody, r.Err
}
func (r *Request) Run(body io.Reader, out interface{}) (statusCode int, responseBody []byte, err error) {
	r.fillOut(out)
	return r.StatusCode, r.ResponseBody, r.Err
}
