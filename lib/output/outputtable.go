package output

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Outputtable is an interface that means the object has a default set of fields and an implementation of prettyprint.PrettyPrinter.
// This means it can be output as a table or list, or in a simpler human-readable format.
type Outputtable interface {
	prettyprint.PrettyPrinter
	DefaultFieldsHaver
}
