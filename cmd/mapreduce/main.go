package main

import (
	"fmt"
	"map-reduce/config"
	"map-reduce/internal"
)

func main() {
	job := &internal.Job{
		InputFile:   "data/input/random_letters.txt",
		NumMappers:  0,
		NumReducers: 0,
	}
	jobConfig := internal.JobConfig{
		ChunkSize:     config.CHUNK_SIZE,
		MaxWorkers:    config.MAX_WORKERS,
		MaxTasks:      config.MAX_TASKS,
		WorkerTimeout: config.WORKER_TIMEOUT,
		RetryLimit:    config.RETRY_LIMIT,
		TempDir:       config.TEMP_DIR,
		OutDir:        config.OUTPUT_DIR,
	}

	job.NumMappers = jobConfig.MapTasks
	job.NumReducers = jobConfig.ReduceTasks

	master := internal.NewMaster(job, jobConfig)
	if master == nil {
		fmt.Println("master failure")
		return
	}
	master.RunMapPhase()
	fmt.Println("map phase completed.")
	master.RunReducePhase()

}
