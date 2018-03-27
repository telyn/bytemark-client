package util_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/util"
)

func TestNumberOrStringInt(t *testing.T) {
	tests := []struct {
		name      string
		nos       util.NumberOrString
		number    int
		shouldErr bool
	}{
		{
			name:   "WithANumber",
			nos:    util.NumberOrString("123"),
			number: 123,
		},
		{
			name:      "WithAString",
			nos:       util.NumberOrString("test"),
			number:    0,
			shouldErr: true,
		},
		{
			name:      "WithEmptyString",
			nos:       util.NumberOrString(""),
			number:    0,
			shouldErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			number, err := test.nos.Int()
			if number != test.number {
				t.Errorf("unexpected number returned: %d", number)
			}
			if test.shouldErr && err == nil {
				t.Error("expected an error")
			}
			if !test.shouldErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNumberOrStringString(t *testing.T) {
	tests := []struct {
		name string
		nos  util.NumberOrString
		str  string
	}{
		{
			name: "WithANumber",
			nos:  util.NumberOrString("123"),
			str:  "123",
		},
		{
			name: "WithAString",
			nos:  util.NumberOrString("test"),
			str:  "test",
		},
		{
			name: "WithEmptyString",
			nos:  util.NumberOrString(""),
			str:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := test.nos.String()
			if str != test.str {
				t.Errorf("unexpected string returned: %s", str)
			}
		})
	}
}

func TestNumberOrStringMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		nos  util.NumberOrString
		json string
	}{
		{
			name: "WithANumber",
			nos:  util.NumberOrString("123"),
			json: `123`,
		},
		{
			name: "WithAString",
			nos:  util.NumberOrString("test"),
			json: `"test"`,
		},
		{
			name: "WithEmptyString",
			nos:  util.NumberOrString(""),
			json: `""`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := test.nos.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			jsonStr := string(jsonData)
			if jsonStr != test.json {
				t.Errorf("unexpected json returned: %s", jsonStr)
			}
		})
	}
}

func TestNumberOrStringUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		nos       util.NumberOrString
		shouldErr bool
	}{
		{
			name: "WithANumber",
			json: `123`,
			nos:  util.NumberOrString("123"),
		},
		{
			name: "WithAString",
			json: `"test"`,
			nos:  util.NumberOrString("test"),
		},
		{
			name: "WithEmptyString",
			json: `""`,
			nos:  util.NumberOrString(""),
		},
		{
			name:      "WithInvalidJSON",
			json:      `[]`,
			shouldErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nos := util.NumberOrString("should be overwritten")
			jsonData := []byte(test.json)
			err := nos.UnmarshalJSON(jsonData)
			if !test.shouldErr && err != nil {
				t.Fatal(err)
			}
			if test.shouldErr {
				if err == nil {
					t.Fatal("expected error")
				} else {
					return
				}
			}
			if nos != test.nos {
				t.Errorf("unexpected value returned: %v", nos)
			}
		})
	}
}
