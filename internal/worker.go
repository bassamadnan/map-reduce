package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"
)

func SerializeKeyValue(kv KeyValue) string {
	keyType := reflect.TypeOf(kv.Key).String()
	valueType := reflect.TypeOf(kv.Value).String()
	valueStr := fmt.Sprintf("%v", kv.Value)
	return fmt.Sprintf("%s\t%s\t%s\t%s", kv.Key, valueStr, keyType, valueType)
}

var globalFileCounter int32 = 0

func (w *Worker) WorkerMapTask(task *Task, config JobConfig) TaskResult {
	fmt.Printf("starting -> worker id: %v, taskid: %v, firstchar: %v, lastchar: %v, sz: %v\n",
		w.ID, task.ID, string(task.Input[0]), string(task.Input[len(task.Input)-1]), len(task.Input))

	mappedData := Map("", task.Input)

	fmt.Printf("ending -> worker id: %v, taskid: %v\n", w.ID, task.ID)

	fileCounter := atomic.AddInt32(&globalFileCounter, 1)
	outputFile := filepath.Join(config.TempDir, fmt.Sprintf("map_%d_%d.txt", task.ID, fileCounter))
	file, err := os.Create(outputFile)
	if err != nil {
		return TaskResult{TaskID: task.ID, Error: err}
	}
	defer file.Close()

	for _, kv := range mappedData {
		serialized := SerializeKeyValue(kv)
		_, err := fmt.Fprintln(file, serialized)
		if err != nil {
			return TaskResult{TaskID: task.ID, Error: err}
		}
	}

	return TaskResult{TaskID: task.ID, Result: outputFile}
}

func (w *Worker) WorkerReduceTask(task *Task, config JobConfig) TaskResult {
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

func (w *Worker) Run(config JobConfig) {
	for task := range w.TaskChannel {
		var result TaskResult
		if task.Type == MapTask {
			result = w.WorkerMapTask(task, config)
		} else {
			result = w.WorkerReduceTask(task, config)
		}
		w.ResultChannel <- result
	}
	fmt.Printf("Worker exitting %v\n", w.ID)
}
