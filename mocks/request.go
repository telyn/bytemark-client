package mocks

import (
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
)

// Request allows bytemark-client developers to stub out requests to bytemark APIs
// in the simple case where a test will only make one API request, Client.MockRequest can be used
// in cases where multiple APIs may be called, use Client.When("BuildRequest").Returns(Request{}) to stub out
// each request individually.
// See the example functions in request_examples_test.go
type Request struct {
	// testingT is an interface that *testing.T implements. See ExampleSingleRequest / ExampleMultiRequest in request_examples_test.go
	T TestingT

	// StatusCode is the status code to be returned by Run / MarshalAndRun
	StatusCode int
	// ResponseBody is the bytes that will be returned from Run/MarshalAndRun.
	// all requests which expect JSON back should use ResponseObject instead
	// since ResponseBody will not be unmarshalled into `out`
	ResponseBody []byte
	// ResponseObject is the object that will be assigned to `out` when
	// MarshalAndRun / Run is called.
	// All requests  which expect a non-JSON response should use ResponseBody
	// instead since ResponseObject cannot be automatically marshalled into
	// the response.
	ResponseObject interface{}
	// Err is an error for Run/MarshalAndRun to return to simulate network or API errors
	Err error

	// requestBody and requestObject are set by Run and MarshalAndRun respectively.
	// make assertions against them with AssertRequestBodyEqual or
	// AssertRequestObjectEqual
	requestBody   io.Reader
	requestObject interface{}
}

func (r *Request) AssertRequestBodyEqual(expected string) {
	if r.requestBody == nil {
		if expected == "" {
			return
		}
		r.T.Fatalf("AssertRequestBodyEqual: Expected was not blank, but requestBody was nil. Was Request.Run actually called?")
	}
	actual, err := ioutil.ReadAll(r.requestBody)
	if err != nil {
		r.T.Fatalf("Couldn't read from requestBody - %s", err)
	}
	if string(actual) != expected {
		r.T.Fatalf("Request body did not equal expected:\nexpected: %q \n  actual: %q", expected, actual)
	}
}

func (r *Request) AssertRequestObjectEqual(expected interface{}) {
	if !reflect.DeepEqual(expected, r.requestObject) {
		r.T.Fatalf("Request body did not equal expected:\nexpected: %#v \n  actual: %#v", expected, r.requestObject)
	}
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
	if out == nil {
		return
	}
	if r.ResponseObject == nil {
		r.T.Errorf("mocks.Request.Run tried to set 'out' object, but the Request.ResponseObject was nil - specify a ResponseObject in your test")
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
	r.requestObject = in
	r.fillOut(out)
	return r.StatusCode, r.ResponseBody, r.Err
}
func (r *Request) Run(body io.Reader, out interface{}) (statusCode int, responseBody []byte, err error) {
	r.fillOut(out)
	r.requestBody = body
	return r.StatusCode, r.ResponseBody, r.Err
}
