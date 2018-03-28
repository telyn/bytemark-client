package admin

import "github.com/BytemarkHosting/bytemark-client/lib/util"

func stringsToNumberOrStrings(in []string) []util.NumberOrString {
	out := make([]util.NumberOrString, len(in))
	for i, str := range in {
		out[i] = util.NumberOrString(str)
	}

	return out
}
