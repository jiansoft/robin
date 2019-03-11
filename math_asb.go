package robin

// Abs Returns the absolute value of a specified int number.
func Abs(a int) int {
	return (a ^ a>>31) - a>>31
}

// AbsForInt64 Returns the absolute value of a specified int64 number.
func AbsForInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
