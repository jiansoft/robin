package robin

import "fmt"

type PendingTask struct {
	identifyId string
	task       Task
	cancelled  bool
}

func (p *PendingTask) init(task Task) *PendingTask {
	p.task = task
	p.cancelled = false
	p.identifyId = fmt.Sprintf("%p-%p", &p, &task)
	return p
}

func NewPendingTask(task Task) *PendingTask {
	return new(PendingTask).init(task)
}

func (p *PendingTask) Dispose() {
	p.cancelled = true
}

func (p *PendingTask) Identify() string {
	return p.identifyId
}

func (p PendingTask) execute() {
	if p.cancelled {
		return
	}
	p.task.Run()
}
