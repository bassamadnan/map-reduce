package internal

import "time"

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
	Inputs map[string][]string
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
	InputPath   string
	OutputPath  string
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
	InputDir      string
	TempDir       string
	OutputDir     string
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
			PendingMapTasks:    config.MapTasks,
			PendingReduceTasks: config.ReduceTasks,
		},
	}

	for i := 0; i < config.MapTasks; i++ {
		master.MapTasks = append(master.MapTasks, MapTask{
			ID: i,
			Status: TaskStatus{
				Worker:    -1,
				Completed: false,
			},
		})
	}

	for i := 0; i < config.ReduceTasks; i++ {
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
