package config

import (
	"path"
	"runtime"
)

const BASE_PATH = "data"

var (
	INPUT_DIR            = path.Join(BASE_PATH, "input")
	TEMP_DIR             = path.Join(BASE_PATH, "temp")
	OUTPUT_DIR           = path.Join(BASE_PATH, "output")
	TOTAL_WORKER_THREADS = runtime.NumCPU()
)

const (
	CHUNK_SIZE     = 1024 * 1024
	MAX_WORKERS    = 10
	MAX_TASKS      = 10
	WORKER_TIMEOUT = 10
	MASTER_PORT    = 8080
	MAP_TASKS      = 4
	REDUCE_TASKS   = 2
	RETRY_LIMIT    = 3
)

type JobConfig struct {
	ChunkSize     int
	MaxWorkers    int
	MaxTasks      int
	WorkerTimeout int
	MasterPort    int
	MapTasks      int
	ReduceTasks   int
	RetryLimit    int
	InputDir      string
	TempDir       string
	OutputDir     string
}
