package output

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type Outputtable interface {
	prettyprint.PrettyPrinter
	DefaultFieldsHaver
}
