package robin

import "testing"

func TestAbs(t *testing.T) {
	type args struct {
		a int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test_Abs_1", args{a: 1}, 1},
		{"Test_Abs_2", args{a: -1}, 1},
		{"Test_Abs_3", args{a: -123}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.a); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
