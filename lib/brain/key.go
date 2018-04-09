package brain

import (
	"fmt"
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Key is an SSH public key which can be output nicely in a table.
type Key struct {
	Key string
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (k Key) DefaultFields(f output.Format) string {
	return "Key"
}

// PrettyPrint writes the key to the given writer
func (k Key) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	_, err := fmt.Fprint(wr, k.Key)
	return err
}

func (k *Key) UnmarshalText(text []byte) error {
	k.Key = string(text)
	return nil
}

// String returns this key as a simple string.
func (k Key) String() string {
	return k.Key
}

type Keys []Key

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (k Keys) DefaultFields(f output.Format) string {
	return "Key"
}

// PrettyPrint outputs the keys, one per line, with no indent.
func (k Keys) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, key := range k {
		_, err = fmt.Fprintln(wr, key.Key)
		if err != nil {
			return
		}
	}
	return
}

func (k *Keys) UnmarshalText(text []byte) error {
	for _, line := range strings.Split(string(text), "\n") {
		*k = append(*k, Key{Key: line})
	}
	return nil
}

func (k Keys) Strings() (strs []string) {
	strs = make([]string, 0)
	for _, key := range k {
		strs = append(strs, key.Key)
	}
	return
}

func (k Keys) MarshalText() ([]byte, error) {
	return []byte(strings.Join(k.Strings(), "\n")), nil
}
