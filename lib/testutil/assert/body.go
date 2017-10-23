package assert

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// BodyString returns a RequestAssertFunc that asserts that the request body
func BodyString(expected string) RequestAssertFunc {
	return Body(func(t *testing.T, testName string, body string) {
		if body != strings.TrimSpace(expected) {
			t.Errorf("%s request body was wrong\nexpected: %q\nactual: %q", testName, expected, body)
		}
	})
}

// BodyUnmarshal returns a RequestAssertFunc that unmarshals the request body
// into object, then runs assertFunc.
// Note that assertFunc does not accept an object - this is because it is intended that you write
// assertFunc as a closure over your object.
func BodyUnmarshal(object interface{}, assertFunc func(*testing.T, string)) RequestAssertFunc {
	return Body(func(t *testing.T, testName string, body string) {
		err := json.Unmarshal([]byte(body), object)
		if err != nil {
			t.Fatalf("%s couldn't unmarshal body: %s", testName, err)
		}

		assertFunc(t, testName)
	})
}

type fakeBody struct {
	buf *bytes.Buffer
}

func (fb *fakeBody) Read(p []byte) (int, error) {
	return fb.buf.Read(p)
}

func (fb *fakeBody) Close() error {
	return nil
}

func BodyUnmarshalEqual(expected map[string]interface{}) RequestAssertFunc {
	body := make(map[string]interface{})
	return BodyUnmarshal(&body, func(t *testing.T, testName string) {
		Equal(t, testName, expected, body)
	})
}

// Body reads the request's body and checks it's the same as expected
func Body(assertFunc func(t *testing.T, testName string, body string)) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("%s couldn't read request body - %s", testName, err)
		}
		_ = r.Body.Close()
		r.Body = &fakeBody{bytes.NewBuffer(body)}
		assertFunc(t, testName, strings.TrimSpace(string(body)))
	}
}

// BodyFormValue returns a RequestAssertFunc that asserts that the body is
// form-encoded and has the key-value pair given.
func BodyFormValue(key, expectedValue string) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Errorf("%s http.Request.ParseForm failed: %s", testName, err)
		}
		URLValue(t, testName, r.Form, key, expectedValue)
	}
}
