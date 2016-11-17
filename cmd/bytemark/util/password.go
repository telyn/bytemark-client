package util

import (
	"math/rand"
	"time"
)

// GeneratePassword generates a random 16-character password made entirely of letters.
func GeneratePassword() (pass string) {
	rand.Seed(time.Now().UnixNano())
	base := byte('A')
	// generate a bitfield of whether-characters-should-be-lowercase
	caps := rand.Uint32()

	for i := uint(0); i < 16; i++ {
		// not lowercase by default
		lc := byte(0)
		// lowercase if the big in position i is 1
		if caps&(0x1<<i) == 0 {
			lc = 0x20 // difference betweeen 'a' and 'A'
		}
		// letter between A and Z, indexed from A
		letter := byte(rand.Intn(25) & 0x1F)
		// add base, lc and letter to get the ascii, and then add that character to pass.
		pass += string(base + lc + letter)

	}

	return pass
}
