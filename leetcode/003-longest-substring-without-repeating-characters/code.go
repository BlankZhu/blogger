func lengthOfLongestSubstring(s string) int {
	dict := make(map[rune]int)
	start := -1
	ret := 0
	for i, v := range s {
		if _, ok := dict[v]; ok {
			if dict[v] > start {
				start = dict[v]
			}
		}
		dict[v] = i
		if ret <= i-start {
			ret = i - start
		}
	}

	return ret
}