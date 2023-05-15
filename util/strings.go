package util

// InArrayString returns true if s is in arr.
func InArrayString(s string, arr []string) bool {
	for _, str := range arr {
		if s == str {
			return true
		}
	}
	return false
}
