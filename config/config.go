package config

import (
	"path"
)

const BASE_PATH = "data"

var (
	INPUT_DIR  = path.Join(BASE_PATH, "input")
	TEMP_DIR   = path.Join(BASE_PATH, "temp")
	OUTPUT_DIR = path.Join(BASE_PATH, "output")
)

const (
	CHUNK_SIZE     = 1024 * 1024
	MAX_WORKERS    = 3
	MAX_TASKS      = 10
	WORKER_TIMEOUT = 10
)
