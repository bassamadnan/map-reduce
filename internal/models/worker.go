package models

type Worker struct {
	ID           int
	Master       *Master // for communication
	AssignedTask *Task
}

type Task interface {
	GetID() int
	GetStatus() TaskStatus
	SetStatus(status TaskStatus)
}

func (m MapTask) GetID() int              { return m.ID }
func (m MapTask) GetStatus() TaskStatus   { return m.Status }
func (m *MapTask) SetStatus(s TaskStatus) { m.Status = s }

func (r ReduceTask) GetID() int              { return r.ID }
func (r ReduceTask) GetStatus() TaskStatus   { return r.Status }
func (r *ReduceTask) SetStatus(s TaskStatus) { r.Status = s }
