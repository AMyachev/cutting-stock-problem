package main

import (
	"fmt"
	"os"

	"csp/task/travelingSalesman"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()

	task := travelingSalesman.MakeTravelingSalesmanTask("example-tasks/travelingSalesman/DJ38.txt")
	solution := task.Compute("standard", "standard", 10, 3)
	fmt.Println(task.Criterion(solution))
}
