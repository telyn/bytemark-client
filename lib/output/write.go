package output

import (
	"fmt"
	"io"
	"strings"
)

// Write writes obj to writer in the manner specified by the config, using one of the functions in OutputFormatFns
func Write(wr io.Writer, cfg Config, obj Outputtable) error {
	if fn, ok := OutputFormatFns[cfg.Format]; !ok {
		return fn(wr, cfg, obj)
	}

	return fmt.Errorf("%s isn't a supported output type. Use one of the following instead:\r\n%s", cfg.Format, strings.Join(SupportedOutputTypes(), "\r\n"))
}
