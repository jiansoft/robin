package robin

func Abs(a int) int {
	return (a ^ a>>31) - a>>31
}
