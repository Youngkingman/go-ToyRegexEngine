package dfa

func compareSet(a map[int]bool, b map[int]bool) bool {
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}
