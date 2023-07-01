package robin

type executor interface {
	executeTasks(t []task)
	executeTask(task)
	executeTasksWithGoroutine(t []task)
	executeTaskWithGoroutine(t task)
}

type defaultExecutor struct {
}

func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (de defaultExecutor) executeTasks(tasks []task) {
	for _, task := range tasks {
		task.execute()
	}
}

func (de defaultExecutor) executeTask(task task) {
	task.execute()
}

func (de defaultExecutor) executeTasksWithGoroutine(tasks []task) {
	for _, task := range tasks {
		de.executeTaskWithGoroutine(task)
	}
}

func (de defaultExecutor) executeTaskWithGoroutine(task task) {
	go task.execute()
}
