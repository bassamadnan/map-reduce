package main

import (
	"fmt"
	"map-reduce/config"
	"map-reduce/internal"
	"map-reduce/pkg"
)

func main() {
	fmt.Println("MapReduce...")
	internal.Process()
	pkg.Utility()
	println(config.INPUT_DIR)
}
