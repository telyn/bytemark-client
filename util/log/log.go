package log

import (
	"fmt"
	"os"
)

const (
	// LvlOutline shows client version & arguments passed in, HTTP URLs & status codes, subprocess calls.
	LvlOutline = 1 + iota
	// LvlArgs shows args going in to library functions
	LvlArgs
	// LvlFlags is used to show/hide the way flags and command line arguments are being messed with in main
	LvlFlags
	// LvlHTTPData is used to show/hide raw HTTP request and response bodies.
	LvlHTTPData
	// LvlMisc is used for any other minutiae
	LvlMisc
)

var DebugLevel int
var LogFile *os.File

func Error(stuff ...interface{}) {
	if len(stuff) == 0 {
		Error("")
	}
	for _, v := range stuff {
		fmt.Fprintln(os.Stderr, v)
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}

}

func Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Log(stuff ...interface{}) {
	if len(stuff) == 0 {
		Log("")
	}
	for _, v := range stuff {
		fmt.Fprintln(os.Stderr, v)
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Output(stuff ...interface{}) {

	if len(stuff) == 0 {
		Output("")
	}
	for _, v := range stuff {
		fmt.Println(v)
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Outputf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Debug(level int, stuff ...interface{}) {
	for _, v := range stuff {
		if level <= DebugLevel {
			fmt.Fprintln(os.Stderr, v)
		}
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Debugf(level int, format string, args ...interface{}) {
	if level <= DebugLevel {
		fmt.Fprintf(os.Stderr, format, args...)
	}
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}
