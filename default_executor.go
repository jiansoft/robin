package robin

type Executor interface {
	ExecuteTasks(t []task)
	ExecuteTasksWithGoroutine(t []task)
}

type defaultExecutor struct {
}

func NewDefaultExecutor() defaultExecutor {
	return defaultExecutor{}
}

func (d defaultExecutor) ExecuteTasks(tasks []task) {
	for _, task := range tasks {
		task.run()
	}
}

func (d defaultExecutor) ExecuteTasksWithGoroutine(tasks []task) {
	for _, task := range tasks {
		go task.run()
	}
}
