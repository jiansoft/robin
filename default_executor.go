package robin

type Executor interface {
	ExecuteTasks(t []Task)
	ExecuteTasksWithGoroutine(t []Task)
}

type defaultExecutor struct {
}

func NewDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (d defaultExecutor) ExecuteTasks(tasks []Task) {
	for _, task := range tasks {
		task.Run()
	}
}

func (d defaultExecutor) ExecuteTasksWithGoroutine(tasks []Task) {
	for _, task := range tasks {
		go task.Run()
	}
}
