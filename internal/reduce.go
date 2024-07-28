package internal

import (
	"fmt"
	"time"
)

func (m *Master) RunReducePhase() {

	// 0. Aggregate data from map phase
	// 1. start Workers
	// 2. assign Tasks
	// 3. wait for completion
	aggregateData, err := m.mergeFiles()
	if err != nil {
		fmt.Println("error aggregating data: ", err)
	}

	m.createReduceTasks(aggregateData)
	m.startWorkers()
	for !m.allReduceTasksCompleted() {
		m.assignTasks()
		select {
		case result := <-m.ResultChannel:
			if result.Error != nil {
				fmt.Println("error : ", result.Error)
			} else {
				for i := range m.Tasks {
					task := &m.Tasks[i]
					if task.ID == result.TaskID && task.Type == ReduceTask {
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
			fmt.Println("Time out reduce phase")

		}
	}
	println("reduce phase complete")
}

func (m *Master) allReduceTasksCompleted() bool {
	for _, task := range m.Tasks {
		if task.Type == ReduceTask && task.Status != TaskCompleted {
			return false
		}
	}
	return true
}

func (m *Master) createReduceTasks(aggregateData map[string][]string) {
	keys := make([]string, 0, len(aggregateData))
	for k := range aggregateData {
		keys = append(keys, k)
	}
	fmt.Println(keys)
	taskLimit := 5 //5 keys per task for reduce func TOOD : set this as config
	numReduceTasks := (len(keys) + taskLimit - 1) / taskLimit
	m.Tasks = make([]Task, numReduceTasks)
	for i := 0; i < numReduceTasks; i++ {
		start := i * taskLimit
		end := min(len(keys), (i+1)*taskLimit)
		taskKeys := keys[start:end]
		taskData := make(map[string][]string)
		for _, key := range taskKeys {
			taskData[key] = aggregateData[key]
		}
		m.Tasks[i] = Task{
			ID:     i,
			Type:   ReduceTask,
			Status: TaskPending,
			Worker: -1,
			Inputs: taskData,
		}
	}
}
