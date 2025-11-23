package util

func ContainsUint64(needle uint64, haystack []uint64) bool {
	var isContains bool = false
	for _, value := range haystack {
		if value == needle {
			isContains = true
		}
	}
	return isContains
}