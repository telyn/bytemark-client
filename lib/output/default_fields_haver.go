package output

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type DefaultFieldsHaver interface {
	DefaultFields(f Format) string
}
