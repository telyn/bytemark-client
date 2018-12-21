package morestrings

import "strings"

// JoinWithSpecialLast joins multiple strings together with a separator, except
// the last two, which are seperated by a different seperator.
func JoinWithSpecialLast(sep string, fin string, strs []string) string {
	// special cases for when there are 0, 1 or 2 strings
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	case 2:
		return strs[0] + fin + strs[1]
	}
	// join with one seperator
	most := strings.Join(strs[0:len(strs)-1], sep)
	// and add the last with the 'fin' seperator.
	return most + fin + strs[len(strs)-1]
}
