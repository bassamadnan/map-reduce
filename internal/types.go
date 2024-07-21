package internal

import "time"

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
)

type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
)

type Task struct {
	ID        int
	Type      TaskType
	Status    TaskStatus
	Worker    int
	Input     string              // For MapTask: input chunk, for ReduceTask: empty
	Inputs    map[string][]string // For ReduceTask: intermediate data, for MapTask: empty
	StartTime time.Time
	EndTime   time.Time
	LastPing  time.Time // remove?
}

type WorkerStatus int

const (
	WorkerIdle WorkerStatus = iota
	WorkerBusy
)

type TaskResult struct {
	TaskID int
	Result interface{}
	Error  error
}

type Job struct {
	MapFunc     MapFunc
	ReduceFunc  ReduceFunc
	InputFile   string
	OutputFile  string
	NumMappers  int
	NumReducers int
}

type JobState struct {
	AvailableWorkers   int
	PendingMapTasks    int
	PendingReduceTasks int
}

type JobConfig struct {
	ChunkSize     int
	MaxWorkers    int
	MaxTasks      int
	WorkerTimeout int
	MasterPort    int
	MapTasks      int
	ReduceTasks   int
	RetryLimit    int
	TempDir       string
}

type Master struct {
	Job           *Job
	Config        JobConfig
	State         JobState
	Tasks         []Task
	Workers       []Worker
	ResultChannel chan TaskResult
}
