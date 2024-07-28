package main

import (
	"map-reduce/config"
	"map-reduce/internal"
)

func initialize() (*internal.Job, internal.JobConfig) {
	job := &internal.Job{
		InputFile: "data/input/random_letters.txt",
		// outdir set by default data/output/
	}

	jobConfig := internal.JobConfig{
		ChunkSize:     config.CHUNK_SIZE,
		MaxWorkers:    config.MAX_WORKERS,
		MaxTasks:      config.MAX_TASKS,
		WorkerTimeout: config.WORKER_TIMEOUT,
		TempDir:       config.TEMP_DIR,
		OutDir:        config.OUTPUT_DIR,
	}
	return job, jobConfig
}
