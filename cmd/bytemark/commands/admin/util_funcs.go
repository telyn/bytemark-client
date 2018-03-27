package admin

import "github.com/BytemarkHosting/bytemark-client/lib/util"

func stringsToJSONNumbers(in []string) (out []util.NumberOrString) {
	out = make([]util.NumberOrString, len(in))
	for i, str := range in {
		out[i] = util.NumberOrString(str)
	}

	return out
}
