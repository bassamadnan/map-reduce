package internal

import (
	"fmt"
	utils "map-reduce/pkg"
	"time"
)

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

func NewMaster(job *Job, config JobConfig) *Master {
	master := &Master{
		Job:    job,
		Config: config,
		State: JobState{
			AvailableWorkers:   config.MaxWorkers,
			PendingMapTasks:    0,
			PendingReduceTasks: config.ReduceTasks,
		},
		Workers:       make([]Worker, config.MaxWorkers),
		Tasks:         make([]Task, 0, config.MapTasks+config.ReduceTasks),
		ResultChannel: make(chan TaskResult, config.MaxWorkers),
	}

	for i := 0; i < config.MaxWorkers; i++ {
		master.Workers[i] = Worker{
			ID:            i,
			Status:        WorkerIdle,
			TaskChannel:   make(chan *Task),
			ResultChannel: master.ResultChannel,
		}
	}

	content, err := utils.ReadInput(job.InputFile)
	if err != nil {
		return nil
	}
	// TODO: split chunks based on workers/map tasks? (paper says user defined)
	chunks := utils.SplitInput(content)

	for i, chunk := range chunks {
		master.Tasks = append(master.Tasks, Task{
			ID:     i,
			Type:   MapTask,
			Status: TaskPending,
			Worker: -1,
			Input:  chunk,
		})
	}
	master.State.PendingMapTasks = len(master.Tasks)

	for i := 0; i < job.NumReducers; i++ {
		master.Tasks = append(master.Tasks, Task{
			ID:     len(chunks) + i,
			Type:   ReduceTask,
			Status: TaskPending,
			Worker: -1,
		})
	}

	return master
}

func (m *Master) getNextTask() *Task {
	for i := range m.Tasks {
		if m.Tasks[i].Status == TaskPending {
			return &m.Tasks[i]
		}
	}
	return nil
}

func (m *Master) assignTasks() {
	// find the first idel worker and assign it first incomplete task
	for i := range m.Workers {
		worker := &m.Workers[i]
		if worker.Status == WorkerIdle {
			// find incomplete task
			task := m.getNextTask()
			if task != nil {
				task.Status = TaskInProgress
				task.Worker = worker.ID
				task.StartTime = time.Now()
				worker.TaskChannel <- task
				worker.Status = WorkerBusy
				if task.Type == MapTask {
					m.State.PendingMapTasks--
					fmt.Println(m.State.PendingMapTasks)
				} else {
					m.State.PendingReduceTasks--
				}
			}
		}
	}
}

func (m *Master) StartWorkers() {
	for i := range m.Workers {
		go m.Workers[i].Run(m.Config)
	}
}

func (m *Master) RunMapPhase() {

	// 1. start Workers
	// 2. assign Tasks
	// 3. collect results (intermediate, to be passed to reduce phase)

	m.StartWorkers()
	for m.State.PendingMapTasks > 0 {
		m.assignTasks()
		select {
		case result := <-m.ResultChannel:
			if result.Error != nil {
				fmt.Printf("task fail --> %v\n", result.TaskID)
			} else {
				for i := range m.Tasks {
					if m.Tasks[i].ID == result.TaskID {
						m.Tasks[i].Status = TaskCompleted
						m.Tasks[i].EndTime = time.Now()
						break
					}
				}
			}
		}
	}
	fmt.Println("REACHED HERE ")
	// run till tasks are complete
	for !m.mapTasksCompleted() {
		time.Sleep(100 * time.Millisecond) // slight delay
	}
}

func (m *Master) mapTasksCompleted() bool {
	for i := range m.Tasks {
		if m.Tasks[i].Type == MapTask && m.Tasks[i].Status != TaskCompleted {
			return false
		}
	}
	return true
}
