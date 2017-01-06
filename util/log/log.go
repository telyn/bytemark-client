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

// DebugLevel determines whether or not debugging output should be output to stderr.
var DebugLevel int

// LogFile is the file which bytemark-client is to log to. This can be nil, in which case it won't. Usually ~/.bytemark/debug.log
var LogFile *os.File

// Error outputs stuff to os.Stderr and LogFile, one thing per line.
func Error(stuff ...interface{}) {
	if len(stuff) == 0 {
		Error("")
	}
	for _, v := range stuff {
		_, _ = fmt.Fprintln(os.Stderr, v)
		if LogFile != nil {
			_, _ = fmt.Fprintln(LogFile, v)
		}
	}

}

// Errorf formats the string and outputs it to Stderr and Logfile.
func Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	if LogFile != nil {
		_, _ = fmt.Fprintf(LogFile, format, args...)
	}
}

// Log outputs stuff to os.Stderr and LogFile, one thing per line.
func Log(stuff ...interface{}) {
	if len(stuff) == 0 {
		Log("")
	}
	for _, v := range stuff {
		_, _ = fmt.Fprintln(os.Stderr, v)
		if LogFile != nil {
			_, _ = fmt.Fprintln(LogFile, v)
		}
	}
}

// Logf formats the string and outputs it to Stderr and Logfile.
func Logf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	if LogFile != nil {
		_, _ = fmt.Fprintf(LogFile, format, args...)
	}
}

// Output outputs stuff to os.Stdout and LogFile, one thing per line.
func Output(stuff ...interface{}) {

	if len(stuff) == 0 {
		Output("")
	}
	for _, v := range stuff {
		fmt.Println(v)
		if LogFile != nil {
			_, _ = fmt.Fprintln(LogFile, v)
		}
	}
}

// Outputf formats the string and outputs it to Stdout and Logfile.
func Outputf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	if LogFile != nil {
		_, _ = fmt.Fprintf(LogFile, format, args...)
	}
}

// Debug outputs stuff to LogFile, and to Stderr if DebugLevel >= level. One thing per line.
func Debug(level int, stuff ...interface{}) {
	for _, v := range stuff {
		if level <= DebugLevel {
			_, _ = fmt.Fprintln(os.Stderr, v)
		}
		if LogFile != nil {
			_, _ = fmt.Fprintln(LogFile, v)
		}
	}
}

// Debugf formats the string and outputs it to LogFile, and to Stderr if DebugLevel >= level.
func Debugf(level int, format string, args ...interface{}) {
	if level <= DebugLevel {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
	}
	if LogFile != nil {
		_, _ = fmt.Fprintf(LogFile, format, args...)
	}
}
