package robin

type executor interface {
	executeTasks(t []Task)
	executeTask(Task)
	executeTasksWithGoroutine(t []Task)
	executeTaskWithGoroutine(t Task)
}

type defaultExecutor struct {
}

func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (de defaultExecutor) executeTasks(tasks []Task) {
	for _, task := range tasks {
		task.execute()
	}
}

func (de defaultExecutor) executeTask(task Task) {
	task.execute()
}

func (de defaultExecutor) executeTasksWithGoroutine(tasks []Task) {
	for _, task := range tasks {
		de.executeTaskWithGoroutine(task)
	}
}

func (de defaultExecutor) executeTaskWithGoroutine(task Task) {
	go task.execute()
}
