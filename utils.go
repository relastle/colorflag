package colorflag

func makeOffsets(strs []string) []int {
	max := 0
	for _, str := range strs {
		if len(str) > max {
			max = len(str)
		}
	}

	res := []int{}
	for _, str := range strs {
		res = append(res, max-len(str))
	}
	return res
}
