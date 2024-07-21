package internal

import (
	"fmt"
	utils "map-reduce/pkg"
	"time"
)

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
		ResultChannel: make(chan TaskResult, config.MaxWorkers),
		Tasks:         make([]Task, 0, config.MapTasks+config.ReduceTasks),
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

	// create map tasks
	for i, chunk := range chunks {
		master.Tasks = append(master.Tasks, Task{
			ID:     i,
			Type:   MapTask,
			Status: TaskPending,
			Worker: -1,
			Input:  chunk,
		})
		fmt.Printf("start %v end  %v size %v\n", string(chunk[0]), string(chunk[len(chunk)-1]), len(chunk))
	}
	master.State.PendingMapTasks = len(chunks)

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
					// fmt.Println(m.State.PendingMapTasks)
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
	for !m.allMapTasksCompleted() {
		m.assignTasks()
		select {
		case result := <-m.ResultChannel:
			if result.Error != nil {
				fmt.Printf("task fail --> %v\n", result.TaskID)
				m.State.PendingMapTasks++
			} else {
				for i := range m.Tasks {
					if m.Tasks[i].ID == result.TaskID {
						m.Tasks[i].Status = TaskCompleted
						m.Tasks[i].EndTime = time.Now()
						m.Workers[m.Tasks[i].Worker].Status = WorkerIdle // reset worker status to idle
						fmt.Printf("completed task %v\n", m.Tasks[i].ID)
						break
					}
				}
			}
		case <-time.After(10000 * time.Millisecond): // timeout
			fmt.Println("Time out")

		}
	}

	// for i, task := range m.Tasks {
	// 	if task.Status != TaskCompleted {
	// 		fmt.Printf("task id: %v, i:%v, status:%v\n", task.ID, i, task.Status)
	// 	}
	// }
	fmt.Println("maphase function exit")
}

func (m *Master) mapTasksCompleted() bool {
	for i := range m.Tasks {
		if m.Tasks[i].Type == MapTask && m.Tasks[i].Status != TaskCompleted {
			return false
		}
	}
	return true
}

func (m *Master) allMapTasksCompleted() bool {
	for _, task := range m.Tasks {
		if task.Type == MapTask && task.Status != TaskCompleted {
			return false
		}
	}
	return true
}
