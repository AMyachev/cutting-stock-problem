package main

import (
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
	task.Compute("standard", 10, 2)
}
