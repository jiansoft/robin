package robin

// Abs returns the absolute value for signed integers and floating-point numbers.
// For signed integer minimum values (for example MinInt64), where the positive
// counterpart cannot be represented in the same type, Abs safely saturates to
// that type's maximum value instead of overflowing.
//
// Abs 回傳有號整數與浮點數的絕對值。
// 對於有號整數最小值（例如 MinInt64），由於同型別無法表示其正值，
// Abs 會安全地飽和為該型別的最大值，避免溢位。
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](v T) T {
	if v < 0 {
		// For signed integer minimum values, -v overflows and remains negative.
		// Use -(v+1) to produce the maximum representable value safely.
		// 有號整數最小值在做 -v 時會溢位且仍為負值。
		// 這裡使用 -(v+1) 安全取得同型別最大值。
		if v == -v {
			return -(v + 1)
		}
		return -v
	}

	// Normalize -0 to +0 for floating-point inputs.
	// 浮點數輸入時，將 -0 正規化為 +0。
	if v == 0 {
		return 0
	}

	return v
}
