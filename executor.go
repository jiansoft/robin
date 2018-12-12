package robin

type executor interface {
	ExecuteTasks(t []Task)
	ExecuteTasksWithGoroutine(t []Task)
}

type defaultExecutor struct {
}

func newDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (d defaultExecutor) ExecuteTasks(tasks []Task) {
	for _, task := range tasks {
		task.run()
	}
}

func (d defaultExecutor) ExecuteTasksWithGoroutine(tasks []Task) {
	for _, task := range tasks {
		go task.run()
	}
}
