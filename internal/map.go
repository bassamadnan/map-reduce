package internal

import (
	"fmt"
	utils "map-reduce/pkg"
	"time"
)

func (m *Master) RunMapPhase() {
	// 0. create map tasks
	// 1. start Workers
	// 2. assign Tasks
	// 3. collect results (intermediate, to be passed to reduce phase)

	m.createMapTasks(m.Job) // TODO : error handling?
	m.startWorkers()
	for !m.allMapTasksCompleted() {
		m.assignTasks()
		select {
		case result := <-m.ResultChannel:
			if result.Error != nil {
				fmt.Println("error : ", result.Error)
			} else {
				for i := range m.Tasks {
					task := &m.Tasks[i]
					if task.ID == result.TaskID {
						task.Status = TaskCompleted
						task.EndTime = time.Now()
						task.Output = string(result.Result)
						worker := &m.Workers[task.Worker]
						worker.Status = WorkerIdle // reset worker status to idle
						fmt.Printf("result %v-> task %v , worker %v\n", result.Result, task.ID, worker.ID)
						break
					}
				}
			}
		case <-time.After(10000 * time.Millisecond): // timeout
			fmt.Println("Time out map phase")

		}
	}
	for i := range m.Workers {
		close(m.Workers[i].TaskChannel)
		m.Workers[i] = Worker{
			ID:            i,
			Status:        WorkerIdle,
			TaskChannel:   make(chan *Task),
			ResultChannel: m.ResultChannel,
		}
	}

}

func (m *Master) allMapTasksCompleted() bool {
	for _, task := range m.Tasks {
		if task.Type == MapTask && task.Status != TaskCompleted {
			return false
		}
	}
	return true
}

func (m *Master) createMapTasks(job *Job) error {
	content, err := utils.ReadInput(job.InputFile)
	if err != nil {
		return err
	}
	chunks := utils.SplitInput(content)
	for i, chunk := range chunks {
		m.Tasks = append(m.Tasks, Task{
			ID:     i,
			Type:   MapTask,
			Status: TaskPending,
			Worker: -1,
			Input:  chunk,
		})
		fmt.Printf("start %v end  %v size %v\n", string(chunk[0]), string(chunk[len(chunk)-1]), len(chunk))
	}
	m.State.PendingMapTasks = len(chunks)
	return nil
}
