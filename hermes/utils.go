package main

func stringFound(s string, strings []string) bool {
	for _, str := range strings {
		if s == str {
			return true
		}
	}

	return false
}
