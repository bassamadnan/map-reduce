package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type MapFunc func(key, value string) []KeyValue          // map (k1,v1) → list(k2,v2)
type ReduceFunc func(key string, values []string) string // reduce (k2,list(v2)) → list(v2)

type KeyValue struct {
	Key   string
	Value string
}

type Worker struct {
	ID            int
	Status        WorkerStatus
	TaskChannel   chan *Task // channel passing tasks from master to worker
	ResultChannel chan TaskResult
}

// count the number of characters
func Map(_, value string) []KeyValue {
	charCount := make(map[rune]int)
	for _, char := range value {
		charCount[char]++
	}

	var result []KeyValue
	for char, count := range charCount {
		result = append(result, KeyValue{
			Key:   string(char),
			Value: strconv.Itoa(count),
		})
	}
	return result
}

func (w *Worker) WorkerMapTask(task *Task, config JobConfig) TaskResult {
	mappedData := Map("", task.Input)

	// add intermediate results to temp directory
	outputFile := filepath.Join(config.TempDir, fmt.Sprintf("w%d.txt", w.ID))
	file, err := os.Create(outputFile)
	if err != nil {
		return TaskResult{TaskID: task.ID, Error: err}
	}
	defer file.Close()

	for _, kv := range mappedData {
		_, err := fmt.Fprintf(file, "%s\t%s\n", kv.Key, kv.Value)
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
			// TODO : reduce task

		}
		w.ResultChannel <- result
	}
}
