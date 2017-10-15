package fiber

type executionState int

// iota 初始化後會自動遞增
const (
	created executionState = iota // value --> 0
	running                       // value --> 1
	stopped                       // value --> 2
)
