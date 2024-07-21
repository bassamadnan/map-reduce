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
	TOTAL_WORKER_THREADS = runtime.NumCPU() // unused?
)

const (
	CHUNK_SIZE     = 1024 * 1024
	MAX_WORKERS    = 3
	MAX_TASKS      = 10
	WORKER_TIMEOUT = 10
	RETRY_LIMIT    = 3 // unused?
)
