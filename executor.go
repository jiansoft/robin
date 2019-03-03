package robin

type executor interface {
	ExecuteTasks(t []Task)
	ExecuteTask(Task)
	ExecuteTasksWithGoroutine(t []Task)
	ExecuteTaskWithGoroutine(t Task)
}

type defaultExecutor struct {
}

func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (d defaultExecutor) ExecuteTasks(tasks []Task) {
	for _, task := range tasks {
		d.ExecuteTask(task)
	}
}

func (d defaultExecutor) ExecuteTask(task Task) {
	task.execute()
}

func (d defaultExecutor) ExecuteTasksWithGoroutine(tasks []Task) {
	for _, task := range tasks {
		d.ExecuteTaskWithGoroutine(task)
	}
}

func (d defaultExecutor) ExecuteTaskWithGoroutine(task Task) {
	go task.execute()
}
