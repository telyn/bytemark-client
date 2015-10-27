package log

import (
	"fmt"
	"os"
)

var DebugLevel int
var LogFile *os.File
var Silent bool

func Error(stuff ...interface{}) {
	if len(stuff) == 0 {
		Error("")
	}
	for _, v := range stuff {
		if !Silent {
			fmt.Fprintln(os.Stderr, v)
		}
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}

}

func Errorf(format string, args ...interface{}) {
	if !Silent {
		fmt.Fprintf(os.Stderr, format, args...)
	}
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Log(stuff ...interface{}) {
	if len(stuff) == 0 {
		Log("")
	}
	for _, v := range stuff {
		if !Silent {
			fmt.Fprintln(os.Stderr, v)
		}
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Logf(format string, args ...interface{}) {
	if !Silent {
		fmt.Fprintf(os.Stderr, format, args...)
	}
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Output(stuff ...interface{}) {

	if len(stuff) == 0 {
		Output("")
	}
	for _, v := range stuff {
		if !Silent {
			fmt.Println(v)
		}
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Outputf(format string, args ...interface{}) {
	if !Silent {
		fmt.Printf(format, args...)
	}
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}

func Debug(level int, stuff ...interface{}) {
	for _, v := range stuff {
		if level <= DebugLevel && !Silent {
			fmt.Fprintln(os.Stderr, v)
		}
		if LogFile != nil {
			fmt.Fprintln(LogFile, v)
		}
	}
}

func Debugf(level int, format string, args ...interface{}) {
	if level <= DebugLevel && !Silent {
		fmt.Fprintf(os.Stderr, format, args...)
	}
	if LogFile != nil {
		fmt.Fprintf(LogFile, format, args...)
	}
}
