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

func (d defaultExecutor) executeTasks(tasks []Task) {
	for _, task := range tasks {
		d.executeTask(task)
	}
}

func (d defaultExecutor) executeTask(task Task) {
	task.execute()
}

func (d defaultExecutor) executeTasksWithGoroutine(tasks []Task) {
	for _, task := range tasks {
		d.executeTaskWithGoroutine(task)
	}
}

func (d defaultExecutor) executeTaskWithGoroutine(task Task) {
	go task.execute()
}
