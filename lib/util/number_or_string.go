package util

import (
	"encoding/json"
	"strconv"
)

// NumberOrString is a string that when marshalled/unmarsalled to/from json
// will be represented as a number if strconv.Atoi believes it is a number,
// or a string otherwise.
type NumberOrString string

// Int returns the NumberOrString as an Int, if possible.
func (nos NumberOrString) Int() (int, error) {
	return strconv.Atoi(string(nos))
}

// String returns the NumberOrString as a string, irrespective of whether it
// can be represented as a number.
func (nos NumberOrString) String() string {
	return string(nos)
}

// MarshalJSON marshals the NumberOrString, representing is as a number where
// appropriate.
func (nos NumberOrString) MarshalJSON() ([]byte, error) {
	number, err := nos.Int()
	if err == nil {
		return json.Marshal(number)
	}
	return json.Marshal(nos.String())
}

// UnmarshalJSON unmarshals a NumberOrString accepting either a json number
// or string.
func (nos *NumberOrString) UnmarshalJSON(data []byte) error {
	var number int
	var str string
	if err := json.Unmarshal(data, &number); err == nil {
		*nos = NumberOrString(strconv.Itoa(number))
		return nil
	}
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*nos = NumberOrString(str)
	return nil
}
