package util

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	is := is.New(t)
	for i := 0; i < 100; i++ {
		pass := GeneratePassword()
		t.Logf("Iteration %d: '%s'", i, pass)
		is.Equal(16, len(pass))
		for x, c := range []byte(pass) {
			if c < 'A' || c > 'z' || (c > 'Z' && c < 'a') {
				t.Logf("character %d ('%c') was not in accepted range.", x, c)
				t.Fail()
			}
		}

	}
}
