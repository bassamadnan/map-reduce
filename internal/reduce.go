package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

// user input
func Reduce(key string, values []string) string {
	count := 0
	for _, v := range values {
		n, err := strconv.Atoi(v)
		if err != nil {
			print(err)
			continue
		}
		count += n
	}
	return fmt.Sprintf("%s\t%d", key, count)
}

func (w *Worker) WorkerReduceTask(task *Task, config JobConfig) TaskResult {
	// fmt.Printf("starting -> worker id: %v, taskid: %v, firstchar: %v, lastchar: %v, sz: %v\n",
	// 	w.ID, task.ID, string(task.Input[0]), string(task.Input[len(task.Input)-1]), len(task.Input))

	results := make([]string, 0)
	taskData := task.Inputs
	for key, values := range taskData {
		result := Reduce(key, values)
		results = append(results, result)
	}

	fmt.Printf("ending -> worker id: %v, taskid: %v\n", w.ID, task.ID)

	fileCounter := atomic.AddInt32(&globalFileCounter, 1)
	outputFile := filepath.Join(config.OutDir, fmt.Sprintf("reduce_%d_%d.txt", task.ID, fileCounter))
	file, err := os.Create(outputFile)
	if err != nil {
		return TaskResult{TaskID: task.ID, Error: err}
	}
	defer file.Close()

	for _, result := range results {
		_, err := fmt.Fprintln(file, result)
		if err != nil {
			return TaskResult{TaskID: task.ID, Error: err}
		}
	}

	return TaskResult{TaskID: task.ID, Result: outputFile}
}
