package output

import (
	"fmt"
	"io"
	"strings"
)

// Write writes obj to writer in the manner specified by the config, using one of the functions in OutputFormatFns
func Write(wr io.Writer, cfg Config, obj Outputtable) error {
	if obj == nil {
		return fmt.Errorf("Object passed to output.Write was nil")
	}
	if fn, ok := FormatFns[cfg.Format]; ok {
		return fn(wr, cfg, obj)
	}
	return fmt.Errorf("%q isn't a supported output type. Use one of the following instead:\r\n%s", cfg.Format, strings.Join(SupportedOutputFormats(), "\r\n"))

}
