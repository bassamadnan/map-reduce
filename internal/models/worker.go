package models

type Worker struct {
	ID           int
	Status       int
	AssignedTask *Task
}
type WorkerStatus int

const (
	WorkerIdle WorkerStatus = iota
	WorkerBusy
	// worker dead?
)

type Task interface {
	GetID() int
	GetStatus() TaskStatus
	SetStatus(status TaskStatus)
	GetType() TaskType
}

type TaskType int

const (
	MapTaskType TaskType = iota
	ReduceTaskType
)

func (m MapTask) GetID() int              { return m.ID }
func (m MapTask) GetStatus() TaskStatus   { return m.Status }
func (m *MapTask) SetStatus(s TaskStatus) { m.Status = s }
func (m MapTask) GetType() TaskType       { return MapTaskType }

func (r ReduceTask) GetType() TaskType       { return ReduceTaskType }
func (r ReduceTask) GetID() int              { return r.ID }
func (r ReduceTask) GetStatus() TaskStatus   { return r.Status }
func (r *ReduceTask) SetStatus(s TaskStatus) { r.Status = s }
