package auth

func factorExists(factors []string, factor string) bool {
	for _, f := range factors {
		if f == factor {
			return true
		}
	}

	return false
}
