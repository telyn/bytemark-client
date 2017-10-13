package testutil

import (
	"runtime"
	"strconv"
	"strings"
)

func funcName(name string) string {
	idx := strings.LastIndex(name, ".") + 1
	return name[idx:]
}

func findTestName() (name string) {
	// it seems unlikely that we'll ever get 10 functions calls away
	maxDepth := 64
	i := 1
	for !strings.HasPrefix(name, "Test") {
		if i > maxDepth {
			return "UNKNOWN TEST"
		}
		callerPC, _, _, _ := runtime.Caller(i)
		name = funcName(runtime.FuncForPC(callerPC).Name())
		i++
	}
	return
}

// TestName returns a good name for the current test.
func Name(testNum int) string {
	if testNum < 0 {
		return findTestName()
	}
	return findTestName() + " " + strconv.Itoa(testNum)
}
