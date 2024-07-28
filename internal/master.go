package internal

import (
	"time"
)

func NewMaster(job *Job, config JobConfig) *Master {
	master := &Master{
		Job:    job,
		Config: config,
		State: JobState{
			PendingMapTasks:    0, // to be determined
			PendingReduceTasks: 0, // to be determined
		},
		Workers:       make([]Worker, config.MaxWorkers),
		ResultChannel: make(chan TaskResult, config.MaxWorkers), // workers give result to this channel in both phases
		Tasks:         make([]Task, 0),                          // array of tasks to be completed
	}

	for i := 0; i < config.MaxWorkers; i++ {
		master.Workers[i] = Worker{
			ID:            i,
			Status:        WorkerIdle,
			TaskChannel:   make(chan *Task),
			ResultChannel: master.ResultChannel,
		}
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
	// find the first ideal worker and assign it first incomplete task
	for i := range m.Workers {
		worker := &m.Workers[i]
		if worker.Status == WorkerIdle {
			// find incomplete task
			task := m.getNextTask()
			if task == nil {
				continue // all tasks assigned/finished
			}
			task.Status = TaskInProgress
			task.Worker = worker.ID
			task.StartTime = time.Now()
			worker.TaskChannel <- task
			worker.Status = WorkerBusy
		}
	}
}

func (m *Master) startWorkers() {
	for i := range m.Workers {
		go m.Workers[i].Run(m.Config)
	}
}
