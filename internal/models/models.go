package internal

import "time"

type Job struct {
	MapFunc     MapFunc
	ReduceFunc  ReduceFunc
	InputPath   string
	OutputPaht  string
	NumMappers  int
	NumReducers int
}

type MapFunc func(key, value string) []KeyValue          // map (k1,v1) → list(k2,v2)
type ReduceFunc func(key string, values []string) string // reduce (k2,list(v2)) → list(v2)

type KeyValue struct {
	Key   string
	Value string
}

type TaskStatus struct {
	Worker    int
	Completed bool
	StartTime time.Time
	EndTime   time.Time
	LastPing  time.Time
}

type MapTask struct {
	ID     int
	Inputs map[string][]string
	Status TaskStatus
}

type ReduceTask struct {
	ID     int
	Inputs map[string][]string
	Status TaskStatus
}
