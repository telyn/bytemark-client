package util

import (
	"math/rand"
	"time"
)

func GeneratePassword() (pass string) {
	rand.Seed(time.Now().UnixNano())
	base := byte('A')
	caps := rand.Uint32()

	for i := uint(0); i < 16; i++ {
		lc := byte(0)
		if caps&(0x1<<i) == 0 {
			lc = 0x20 // 'a' - 'A'
		}
		pass += string(base + lc + byte(rand.Intn(25)&31))

	}

	return pass
}
