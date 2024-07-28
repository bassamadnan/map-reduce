package main

import (
	"map-reduce/internal"
)

func main() {

	// get configuration
	job, jobConfig := initialize()

	// create master based on config
	master := internal.NewMaster(job, jobConfig)

	// run map phase
	master.RunMapPhase()

	// run reduce phase
	master.RunReducePhase()

}
