package flags

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//FileFlag implements the flag.Value (aka urfave/cli.Generic) interface
// to provide a flag value type that reads its effective value from the file named as its input.
type FileFlag struct {
	FileName string
	Value    string
}

func getPath(name string) string {
	if len(name) > 0 && name[0] == '~' {
		home := os.Getenv("HOME")
		return filepath.Join(home, name[1:])
	}
	return name
}

// Set sets the value of FileFlag given the filename as an argument. This reads in the file synchronously.
func (f *FileFlag) Set(name string) error {
	f.FileName = getPath(name)
	res, err := ioutil.ReadFile(f.FileName)
	f.Value = string(res)
	return err
}

func (f *FileFlag) String() string {
	return f.Value
}
