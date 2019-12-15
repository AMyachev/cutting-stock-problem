package main

import (
	"fmt"
	"os"

	"csp/task/resourceAllocation"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()
	task := resourceAllocation.MakeResourceAllocationTaskFromFile("example-tasks/resourceAllocation/task_4_01_n2_m2_T2.txt")
	maxFlow := task.Compute()
	fmt.Println(maxFlow)
}
