package internal

import "time"

type MapFunc func(key, value string) []KeyValue          // map (k1,v1) → list(k2,v2)
type ReduceFunc func(key string, values []string) string // reduce (k2,list(v2)) → list(v2)

type KeyValue struct {
	Key   string
	Value interface{}
}

type Worker struct {
	ID            int
	Status        WorkerStatus
	TaskChannel   chan *Task
	ResultChannel chan TaskResult
}

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
	Output    string
}

type WorkerStatus int

const (
	WorkerIdle WorkerStatus = iota
	WorkerBusy
	WorkerFailure
)

type TaskResult struct {
	TaskID int
	Result string
	Error  error
}

type Job struct {
	MapFunc    MapFunc
	ReduceFunc ReduceFunc
	InputFile  string
	OutputFile string
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
	OutDir        string
}

type Master struct {
	Job           *Job
	Config        JobConfig
	State         JobState
	Tasks         []Task
	Workers       []Worker
	ResultChannel chan TaskResult
}
