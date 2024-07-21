package main

import (
	"fmt"
	"map-reduce/config"
	"map-reduce/internal"
)

func main() {

	// master := models.NewMaster(job *models.Job, config models.JobConfig)
	// for worker in master.workers:
	// 	go workerRoutine(worker)

	// while master.maptasks > 0: assign tasks and collect results
	job := &internal.Job{
		MapFunc:     internal.Map,
		InputFile:   "data/input/random_letters.txt",
		OutputFile:  "data/output/result.txt",
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

}
