package brain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalUser(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected Keys
	}{
		{
			name: "oneline",
			in:   `{"authorized_keys":"hello"}`,
			expected: Keys{
				Key{Key: "hello"},
			},
		}, {
			name: "eightlines",
			in:   `{"authorized_keys":"hello\nto\nall\nmy\nmany\nbeautiful\nfriends\nworldwide"}`,
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
			var user User
			err := json.Unmarshal([]byte(test.in), &user)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(test.expected, user.AuthorizedKeys) {
				t.Errorf("expected %#v\ngot      %#v", test.expected, user.AuthorizedKeys)
			}
		})
	}
}
