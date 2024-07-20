package models

import (
	utils "map-reduce/pkg"
	"time"
)

type MapFunc func(key, value string) []KeyValue
type ReduceFunc func(key string, values []string) string

type KeyValue struct {
	Key   string
	Value string
}

type TaskStatus struct {
	Worker    int
	Completed bool
	StartTime time.Time
	EndTime   time.Time
	LastPing  time.Time
}

type MapTask struct {
	ID     int
	Input  string
	Status TaskStatus
}

type ReduceTask struct {
	ID     int
	Inputs map[string][]string
	Status TaskStatus
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
	TempDir       string // remove ?
}

type Master struct {
	Job         *Job
	Config      JobConfig
	State       JobState
	MapTasks    []MapTask
	ReduceTasks []ReduceTask
}

func SpawnMaster(job *Job, config JobConfig) *Master {
	master := &Master{
		Job:    job,
		Config: config,
		State: JobState{
			AvailableWorkers:   config.MaxWorkers,
			PendingMapTasks:    0,
			PendingReduceTasks: config.ReduceTasks,
		},
	}

	content, err := utils.ReadInput(job.InputFile)
	if err != nil {
		return nil
	}
	chunks := utils.SplitInput(content)

	for i, chunk := range chunks {
		master.MapTasks = append(master.MapTasks, MapTask{
			ID:    i,
			Input: chunk,
			Status: TaskStatus{
				Worker:    -1,
				Completed: false,
			},
		})
	}
	master.State.PendingMapTasks = len(master.MapTasks)

	for i := 0; i < job.NumReducers; i++ {
		master.ReduceTasks = append(master.ReduceTasks, ReduceTask{
			ID: i,
			Status: TaskStatus{
				Worker:    -1,
				Completed: false,
			},
		})
	}

	return master
}

func (m *Master) AssignTasks() Task {
	for i := range m.MapTasks {
		if m.MapTasks[i].Status.Completed {
			return &m.MapTasks[i]
		}
	}
	for i := range m.ReduceTasks {
		if m.ReduceTasks[i].Status.Completed {
			return &m.ReduceTasks[i]
		}
	}
	return nil
}

func updateTaskStatus(t *Task, s int) {

}
