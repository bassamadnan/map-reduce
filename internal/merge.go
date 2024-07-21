package internal

import "fmt"

func (m *Master) mergeFiles() {
	fmt.Println("starting merging ..")
	for _, task := range m.Tasks {
		println(task.Output)
	}
}
