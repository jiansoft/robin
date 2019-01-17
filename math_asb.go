package robin

func Abs(a int) int {
	return (a ^ a>>31) - a>>31
}

func AbsForInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
