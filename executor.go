package robin

// executor defines internal task execution strategies for fiber implementations.
// executor 定義 fiber 實作所使用的內部任務執行策略。
type executor interface {
	executeTasks(t []task)
	executeTask(task)
	executeTasksWithGoroutine(t []task)
	executeTaskWithGoroutine(t task)
}

// defaultExecutor is the standard executor used by both single and multi fibers.
// defaultExecutor 是單/多 fiber 共用的標準執行器。
type defaultExecutor struct {
}

// newDefaultExecutor returns a value-type executor with no runtime state.
// newDefaultExecutor 會回傳不含執行期狀態的值型別執行器。
func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

// executeTasks executes tasks sequentially in the caller goroutine.
// executeTasks 會在呼叫端 goroutine 內依序執行所有任務。
func (de defaultExecutor) executeTasks(tasks []task) {
	for _, task := range tasks {
		task.execute()
	}
}

// executeTask executes one task synchronously.
// executeTask 以同步方式執行單一任務。
func (de defaultExecutor) executeTask(task task) {
	task.execute()
}

// executeTasksWithGoroutine executes each task in a separate goroutine.
// executeTasksWithGoroutine 會讓每個任務各自以 goroutine 併發執行。
func (de defaultExecutor) executeTasksWithGoroutine(tasks []task) {
	for _, task := range tasks {
		de.executeTaskWithGoroutine(task)
	}
}

// executeTaskWithGoroutine schedules one task on a new goroutine.
// executeTaskWithGoroutine 會以新的 goroutine 排程執行單一任務。
func (de defaultExecutor) executeTaskWithGoroutine(task task) {
	go task.execute()
}
