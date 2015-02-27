package cmd

func FirstNotEmpty(choices ...string) string {
	for _, choice := range choices {
		if choice != "" {
			return choice

		}
	}
	return ""
}
