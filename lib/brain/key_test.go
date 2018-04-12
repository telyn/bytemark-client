package brain

import (
	"reflect"
	"testing"
)

func TestUnmarshalKeys(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected Keys
	}{
		{
			name: "oneline",
			in:   "hello",
			expected: Keys{
				Key{Key: "hello"},
			},
		}, {
			name: "eightlines",
			in:   "hello\nto\nall\nmy\nmany\nbeautiful\nfriends\nworldwide",
			expected: Keys{
				Key{Key: "hello"},
				Key{Key: "to"},
				Key{Key: "all"},
				Key{Key: "my"},
				Key{Key: "many"},
				Key{Key: "beautiful"},
				Key{Key: "friends"},
				Key{Key: "worldwide"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var keys Keys
			err := keys.UnmarshalText([]byte(test.in))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(test.expected, keys) {
				t.Errorf("expected %#v\ngot      %#v", test.expected, keys)
			}
		})
	}
}
