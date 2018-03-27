package admin

import "encoding/json"

func stringsToJsonNumbers(in []string) (out []json.Number) {
	out = make([]json.Number, len(in))
	for i, str := range in {
		out[i] = json.Number(str)
	}

	return out
}
