package robin

import (
	"math"
	"testing"
)

// TestAbsSignedIntegers verifies Abs() behavior across signed integer types,
// including minimum-value overflow boundaries.
// TestAbsSignedIntegers 驗證 Abs() 在各有號整數型別的行為，
// 包含最小值溢位邊界。
func TestAbsSignedIntegers(t *testing.T) {
	t.Run("int_basic_values", func(t *testing.T) {
		if got := Abs(123); got != 123 {
			t.Fatalf("Abs(123) = %d, want %d", got, 123)
		}
		if got := Abs(-123); got != 123 {
			t.Fatalf("Abs(-123) = %d, want %d", got, 123)
		}
	})

	t.Run("int_min_saturates_to_max", func(t *testing.T) {
		maxInt := int(^uint(0) >> 1)
		minInt := -maxInt - 1
		if got := Abs(minInt); got != maxInt {
			t.Fatalf("Abs(minInt) = %d, want %d", got, maxInt)
		}
	})

	t.Run("int8_min_saturates_to_max", func(t *testing.T) {
		max := int8(^uint8(0) >> 1)
		min := -max - 1
		if got := Abs(min); got != max {
			t.Fatalf("Abs(int8 min) = %d, want %d", got, max)
		}
	})

	t.Run("int16_min_saturates_to_max", func(t *testing.T) {
		max := int16(^uint16(0) >> 1)
		min := -max - 1
		if got := Abs(min); got != max {
			t.Fatalf("Abs(int16 min) = %d, want %d", got, max)
		}
	})

	t.Run("int32_min_saturates_to_max", func(t *testing.T) {
		max := int32(^uint32(0) >> 1)
		min := -max - 1
		if got := Abs(min); got != max {
			t.Fatalf("Abs(int32 min) = %d, want %d", got, max)
		}
	})

	t.Run("int64_min_saturates_to_max", func(t *testing.T) {
		max := int64(^uint64(0) >> 1)
		min := -max - 1
		if got := Abs(min); got != max {
			t.Fatalf("Abs(int64 min) = %d, want %d", got, max)
		}
	})
}

// TestAbsFloats verifies Abs() behavior for floating-point values,
// including negative zero, infinity, and NaN.
// TestAbsFloats 驗證 Abs() 在浮點數上的行為，
// 包含負零、無限大與 NaN。
func TestAbsFloats(t *testing.T) {
	t.Run("float32_negative_value", func(t *testing.T) {
		if got := Abs(float32(-3.5)); got != float32(3.5) {
			t.Fatalf("Abs(float32(-3.5)) = %v, want %v", got, float32(3.5))
		}
	})

	t.Run("float64_negative_value", func(t *testing.T) {
		if got := Abs(-3.5); got != 3.5 {
			t.Fatalf("Abs(-3.5) = %v, want %v", got, 3.5)
		}
	})

	t.Run("float32_negative_zero_normalizes_to_positive_zero", func(t *testing.T) {
		negZero := float32(math.Copysign(0, -1))
		got := Abs(negZero)
		if got != 0 {
			t.Fatalf("Abs(float32(-0)) = %v, want 0", got)
		}
		if math.Signbit(float64(got)) {
			t.Fatalf("Abs(float32(-0)) keeps negative sign bit")
		}
	})

	t.Run("float64_negative_zero_normalizes_to_positive_zero", func(t *testing.T) {
		negZero := math.Copysign(0, -1)
		got := Abs(negZero)
		if got != 0 {
			t.Fatalf("Abs(float64(-0)) = %v, want 0", got)
		}
		if math.Signbit(got) {
			t.Fatalf("Abs(float64(-0)) keeps negative sign bit")
		}
	})

	t.Run("float64_negative_infinity", func(t *testing.T) {
		if got := Abs(math.Inf(-1)); !math.IsInf(got, 1) {
			t.Fatalf("Abs(-Inf) = %v, want +Inf", got)
		}
	})

	t.Run("float64_nan_passthrough", func(t *testing.T) {
		if got := Abs(math.NaN()); !math.IsNaN(got) {
			t.Fatalf("Abs(NaN) = %v, want NaN", got)
		}
	})

	t.Run("float32_nan_passthrough", func(t *testing.T) {
		input := float32(math.NaN())
		if got := Abs(input); !math.IsNaN(float64(got)) {
			t.Fatalf("Abs(float32 NaN) = %v, want NaN", got)
		}
	})
}
