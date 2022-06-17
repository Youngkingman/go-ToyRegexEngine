package dfa

func compareSet(a map[int]bool, b map[int]bool) bool {
	for k := range a {
		_, has := b[k]
		if !has {
			return false
		}
	}
	return true
}
