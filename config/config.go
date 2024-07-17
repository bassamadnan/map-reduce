package config

import "path"

const BASE_PATH = "data"

var (
	INPUT_DIR  = path.Join(BASE_PATH, "input")
	TEMP_DIR   = path.Join(BASE_PATH, "temp")
	OUTPUT_DIR = path.Join(BASE_PATH, "output")
)
