package output

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type DefaultFieldsHaver interface {
	prettyprint.PrettyPrinter
	DefaultFields(f Format) string
}
