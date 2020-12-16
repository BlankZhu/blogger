func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, v := range nums {
		sub := target - v
		if val, ok := numMap[sub]; ok {
			ret := make([]int, 2)
			ret[0] = i
			ret[1] = val
			return ret
		}
		numMap[v] = i
	}
	return make([]int, 0)
}