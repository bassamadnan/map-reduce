package internal

/*
For each map task and reduce task, it stores the state (idle, in-progress, or completed), and
the identity of the worker machine (for non-idle tasks).
*/

type Master struct {
	Job         *Job
	mapTasks    []MapTask
	reduceTasks []ReduceFunc
}
