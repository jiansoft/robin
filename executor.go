package robin

// executor defines task execution behaviors for serial and goroutine modes.
type executor interface {
	executeTasks(t []task)
	executeTask(task)
	executeTasksWithGoroutine(t []task)
	executeTaskWithGoroutine(t task)
}

// defaultExecutor executes tasks using direct calls or goroutines.
type defaultExecutor struct {
}

// newDefaultExecutor creates a defaultExecutor instance.
func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

// executeTasks executes all tasks sequentially.
func (de defaultExecutor) executeTasks(tasks []task) {
	for _, task := range tasks {
		task.execute()
	}
}

// executeTask executes a single task sequentially.
func (de defaultExecutor) executeTask(task task) {
	task.execute()
}

// executeTasksWithGoroutine executes all tasks concurrently with goroutines.
func (de defaultExecutor) executeTasksWithGoroutine(tasks []task) {
	for _, task := range tasks {
		de.executeTaskWithGoroutine(task)
	}
}

// executeTaskWithGoroutine executes a single task in a goroutine.
func (de defaultExecutor) executeTaskWithGoroutine(task task) {
	go task.execute()
}
