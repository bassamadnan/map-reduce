package models

type MapFunc func(key, value string) []KeyValue
type ReduceFunc func(key string, values []string) string

type KeyValue struct {
	Key   string
	Value string
}
