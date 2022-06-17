package mdfa

import "sort"

func compareTwoGamma(a []map[int]bool, b []map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return len(a[i]) > len(a[j])
	})
	sort.Slice(b, func(i, j int) bool {
		return len(b[i]) > len(b[j])
	})
	for i, mp := range a {
		if !compareSet(mp, b[i]) {
			return false
		}
	}
	return true
}

func compareSet(a map[int]bool, b map[int]bool) bool {
	for k := range a {
		_, has := b[k]
		if !has {
			return false
		}
	}
	return true
}
